package binding

import (
	"fmt"
	"path/filepath"

	"k8s.io/klog"

	bindingApis "github.com/redhat-developer/service-binding-operator/apis"
	bindingApi "github.com/redhat-developer/service-binding-operator/apis/binding/v1alpha1"
	specApi "github.com/redhat-developer/service-binding-operator/apis/spec/v1alpha3"

	"github.com/redhat-developer/odo/pkg/project"

	"github.com/devfile/library/pkg/devfile/parser"
	devfilefs "github.com/devfile/library/pkg/testingutil/filesystem"
	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	devfilev1alpha2 "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	parsercommon "github.com/devfile/library/pkg/devfile/parser/data/v2/common"

	"github.com/redhat-developer/odo/pkg/api"
	"github.com/redhat-developer/odo/pkg/binding/asker"
	backendpkg "github.com/redhat-developer/odo/pkg/binding/backend"
	"github.com/redhat-developer/odo/pkg/kclient"
	"github.com/redhat-developer/odo/pkg/libdevfile"
)

type BindingClient struct {
	// Backends
	flagsBackend       *backendpkg.FlagsBackend
	interactiveBackend *backendpkg.InteractiveBackend

	// Clients
	kubernetesClient kclient.ClientInterface
}

var _ Client = (*BindingClient)(nil)

func NewBindingClient(projectClient project.Client, kubernetesClient kclient.ClientInterface) *BindingClient {
	// We create the asker client and the backends here and not at the CLI level, as we want to hide these details to the CLI
	askerClient := asker.NewSurveyAsker()
	return &BindingClient{
		flagsBackend:       backendpkg.NewFlagsBackend(),
		interactiveBackend: backendpkg.NewInteractiveBackend(askerClient, projectClient, kubernetesClient),
		kubernetesClient:   kubernetesClient,
	}
}

// GetFlags gets the flag specific to add binding operation so that it can correctly decide on the backend to be used
// It ignores all the flags except the ones specific to add binding operation, for e.g. verbosity flag
func (o *BindingClient) GetFlags(flags map[string]string) map[string]string {
	bindingFlags := map[string]string{}
	for flag, value := range flags {
		if flag == backendpkg.FLAG_NAME ||
			flag == backendpkg.FLAG_WORKLOAD ||
			flag == backendpkg.FLAG_SERVICE_NAMESPACE ||
			flag == backendpkg.FLAG_SERVICE ||
			flag == backendpkg.FLAG_BIND_AS_FILES ||
			flag == backendpkg.FLAG_NAMING_STRATEGY {
			bindingFlags[flag] = value
		}
	}
	return bindingFlags
}

func (o *BindingClient) GetServiceInstances(namespace string) (map[string]unstructured.Unstructured, error) {
	err := o.checkServiceBindingOperatorInstalled()
	if err != nil {
		return nil, err
	}

	// Get the BindableKinds/bindable-kinds object
	bindableKind, err := o.kubernetesClient.GetBindableKinds()
	if err != nil {
		return nil, err
	}

	// get a list of restMappings of all the GVKs present in bindableKind's Status
	bindableKindRestMappings, err := o.kubernetesClient.GetBindableKindStatusRestMapping(bindableKind.Status)
	if err != nil {
		return nil, err
	}

	var bindableObjectMap = map[string]unstructured.Unstructured{}
	for _, restMapping := range bindableKindRestMappings {
		// TODO: Debug into why List returns all the versions instead of the GVR version
		// List all the instances of the restMapping object
		resources, err := o.kubernetesClient.ListDynamicResources(namespace, restMapping.Resource, "")
		if err != nil {
			if kerrors.IsNotFound(err) || kerrors.IsForbidden(err) {
				// Assume the namespace is deleted or being terminated, hence user can't list its resources
				klog.V(3).Infoln(err)
				continue
			}
			return nil, err
		}

		for _, item := range resources.Items {
			// format: `<name> (<kind>.<group>)`
			serviceName := fmt.Sprintf("%s (%s.%s)", item.GetName(), item.GetKind(), item.GroupVersionKind().Group)
			bindableObjectMap[serviceName] = item
		}

	}

	return bindableObjectMap, nil
}

// GetBindingsFromDevfile returns all ServiceBinding resources declared as Kubernertes component from a Devfile
// from group binding.operators.coreos.com/v1alpha1 or servicebinding.io/v1alpha3
func (o *BindingClient) GetBindingsFromDevfile(devfileObj parser.DevfileObj, context string) ([]api.ServiceBinding, error) {
	result := []api.ServiceBinding{}
	kubeComponents, err := devfileObj.Data.GetComponents(parsercommon.DevfileOptions{
		ComponentOptions: parsercommon.ComponentOptions{
			ComponentType: devfilev1alpha2.KubernetesComponentType,
		},
	})
	if err != nil {
		return nil, err
	}

	for _, component := range kubeComponents {
		strCRD, err := libdevfile.GetK8sManifestsWithVariablesSubstituted(devfileObj, component.Name, context, devfilefs.DefaultFs{})
		if err != nil {
			return nil, err
		}

		u := unstructured.Unstructured{}
		if err := yaml.Unmarshal([]byte(strCRD), &u.Object); err != nil {
			return nil, err
		}

		switch u.GetObjectKind().GroupVersionKind() {
		case bindingApi.GroupVersionKind:

			var sbo bindingApi.ServiceBinding
			err := kclient.ConvertUnstructuredToResource(u, &sbo)
			if err != nil {
				return nil, err
			}

			sb, err := o.kubernetesClient.APIServiceBindingFromBinding(sbo)
			if err != nil {
				return nil, err
			}
			sb.Status, err = o.getStatusFromBinding(sb.Name)
			if err != nil {
				return nil, err
			}

			result = append(result, sb)

		case specApi.GroupVersion.WithKind("ServiceBinding"):

			var sbc specApi.ServiceBinding
			err := kclient.ConvertUnstructuredToResource(u, &sbc)
			if err != nil {
				return nil, err
			}

			sb := o.kubernetesClient.APIServiceBindingFromSpec(sbc)
			sb.Status, err = o.getStatusFromSpec(sb.Name)
			if err != nil {
				return nil, err
			}

			result = append(result, sb)

		}
	}
	return result, nil
}

// GetBindingFromCluster returns the ServiceBinding resource with the given name
// from the cluster, from group binding.operators.coreos.com/v1alpha1 or servicebinding.io/v1alpha3
func (o *BindingClient) GetBindingFromCluster(name string) (api.ServiceBinding, error) {

	bindingSB, err := o.kubernetesClient.GetBindingServiceBinding(name)
	if err == nil {
		var sb api.ServiceBinding
		sb, err = o.kubernetesClient.APIServiceBindingFromBinding(bindingSB)
		if err != nil {
			return api.ServiceBinding{}, err
		}
		sb.Status, err = o.getStatusFromBinding(bindingSB.Name)
		if err != nil {
			return api.ServiceBinding{}, err
		}
		return sb, nil
	}
	if err != nil && !kerrors.IsNotFound(err) {
		return api.ServiceBinding{}, err
	}

	specSB, err := o.kubernetesClient.GetSpecServiceBinding(name)
	if err == nil {
		sb := o.kubernetesClient.APIServiceBindingFromSpec(specSB)
		sb.Status, err = o.getStatusFromSpec(specSB.Name)
		if err != nil {
			return api.ServiceBinding{}, err
		}
		return sb, nil
	}

	// In case of notFound error, this time we return the error
	if kerrors.IsNotFound(err) {
		return api.ServiceBinding{}, fmt.Errorf("ServiceBinding %q not found", name)
	}
	return api.ServiceBinding{}, err
}

// getStatusFromBinding returns status information from a ServiceBinding in the cluster
// from group binding.operators.coreos.com/v1alpha1
func (o *BindingClient) getStatusFromBinding(name string) (*api.ServiceBindingStatus, error) {
	bindingSB, err := o.kubernetesClient.GetBindingServiceBinding(name)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	if injected := meta.IsStatusConditionTrue(bindingSB.Status.Conditions, bindingApis.InjectionReady); !injected {
		return nil, nil
	}

	secretName := bindingSB.Status.Secret
	secret, err := o.kubernetesClient.GetSecret(secretName, o.kubernetesClient.GetCurrentNamespace())
	if err != nil {
		return nil, err
	}

	bindings := make([]string, 0, len(secret.Data))
	if bindingSB.Spec.BindAsFiles {
		for k := range secret.Data {
			bindingName := filepath.ToSlash(filepath.Join("${SERVICE_BINDING_ROOT}", name, k))
			bindings = append(bindings, bindingName)
		}
		return &api.ServiceBindingStatus{
			BindingFiles: bindings,
		}, nil
	}

	for k := range secret.Data {
		bindings = append(bindings, k)
	}
	return &api.ServiceBindingStatus{
		BindingEnvVars: bindings,
	}, nil
}

// getStatusFromSpec returns status information from a ServiceBinding in the cluster
// from group servicebinding.io/v1alpha3
func (o *BindingClient) getStatusFromSpec(name string) (*api.ServiceBindingStatus, error) {
	specSB, err := o.kubernetesClient.GetSpecServiceBinding(name)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	if injected := meta.IsStatusConditionTrue(specSB.Status.Conditions, bindingApis.InjectionReady); !injected {
		return nil, nil
	}

	if specSB.Status.Binding == nil {
		return nil, nil
	}
	secretName := specSB.Status.Binding.Name
	secret, err := o.kubernetesClient.GetSecret(secretName, o.kubernetesClient.GetCurrentNamespace())
	if err != nil {
		return nil, err
	}
	bindingFiles := make([]string, 0, len(secret.Data))
	bindingEnvVars := make([]string, 0, len(specSB.Spec.Env))
	for k := range secret.Data {
		bindingName := filepath.ToSlash(filepath.Join("${SERVICE_BINDING_ROOT}", name, k))
		bindingFiles = append(bindingFiles, bindingName)
	}
	for _, env := range specSB.Spec.Env {
		bindingEnvVars = append(bindingEnvVars, env.Name)
	}
	return &api.ServiceBindingStatus{
		BindingFiles:   bindingFiles,
		BindingEnvVars: bindingEnvVars,
	}, nil
}

func (o *BindingClient) checkServiceBindingOperatorInstalled() error {
	isServiceBindingInstalled, err := o.kubernetesClient.IsServiceBindingSupported()
	if err != nil {
		return err
	}
	if !isServiceBindingInstalled {
		//revive:disable:error-strings This is a top-level error message displayed as is to the end user
		return fmt.Errorf("Service Binding Operator is not installed on the cluster, please ensure it is installed before proceeding. " +
			"See installation instructions: https://odo.dev/docs/command-reference/add-binding#installing-the-service-binding-operator")
		//revive:enable:error-strings
	}
	return nil
}

func (o *BindingClient) CheckServiceBindingsInjectionDone(componentName string, appName string) (bool, error) {

	deployment, err := o.kubernetesClient.GetOneDeployment(componentName, appName, true)
	if err != nil {
		// If not deployment yet => all bindings are done
		if _, ok := err.(*kclient.DeploymentNotFoundError); ok {
			return true, nil
		}
		return false, err
	}
	deploymentName := deployment.GetName()

	specList, bindingList, err := o.kubernetesClient.ListServiceBindingsFromAllGroups()
	if err != nil {
		// If ServiceBinding kind is not registered => all bindings are done
		if runtime.IsNotRegisteredError(err) {
			return true, nil
		}
		return false, err
	}

	for _, binding := range bindingList {
		app := binding.Spec.Application
		if app.Group != appsv1.SchemeGroupVersion.Group ||
			app.Version != appsv1.SchemeGroupVersion.Version ||
			(app.Kind != "Deployment" && app.Resource != "deployments") {
			continue
		}
		if app.Name != deploymentName {
			continue
		}
		if injected := meta.IsStatusConditionTrue(binding.Status.Conditions, bindingApis.InjectionReady); !injected {
			return false, nil
		}
	}

	for _, binding := range specList {
		app := binding.Spec.Workload
		if app.APIVersion != appsv1.SchemeGroupVersion.String() ||
			app.Kind != "Deployment" {
			continue
		}
		if app.Name != deploymentName {
			continue
		}
		if injected := meta.IsStatusConditionTrue(binding.Status.Conditions, bindingApis.InjectionReady); !injected {
			return false, nil
		}
	}

	return true, nil
}
