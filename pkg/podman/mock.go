// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/podman/interface.go

// Package podman is a generated GoMock package.
package podman

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	api "github.com/redhat-developer/odo/pkg/api"
	v1 "k8s.io/api/core/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// CleanupPodResources mocks base method.
func (m *MockClient) CleanupPodResources(pod *v1.Pod) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CleanupPodResources", pod)
	ret0, _ := ret[0].(error)
	return ret0
}

// CleanupPodResources indicates an expected call of CleanupPodResources.
func (mr *MockClientMockRecorder) CleanupPodResources(pod interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanupPodResources", reflect.TypeOf((*MockClient)(nil).CleanupPodResources), pod)
}

// ExecCMDInContainer mocks base method.
func (m *MockClient) ExecCMDInContainer(containerName, podName string, cmd []string, stdout, stderr io.Writer, stdin io.Reader, tty bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecCMDInContainer", containerName, podName, cmd, stdout, stderr, stdin, tty)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExecCMDInContainer indicates an expected call of ExecCMDInContainer.
func (mr *MockClientMockRecorder) ExecCMDInContainer(containerName, podName, cmd, stdout, stderr, stdin, tty interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecCMDInContainer", reflect.TypeOf((*MockClient)(nil).ExecCMDInContainer), containerName, podName, cmd, stdout, stderr, stdin, tty)
}

// GetAllPodsInNamespaceMatchingSelector mocks base method.
func (m *MockClient) GetAllPodsInNamespaceMatchingSelector(selector, ns string) (*v1.PodList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllPodsInNamespaceMatchingSelector", selector, ns)
	ret0, _ := ret[0].(*v1.PodList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllPodsInNamespaceMatchingSelector indicates an expected call of GetAllPodsInNamespaceMatchingSelector.
func (mr *MockClientMockRecorder) GetAllPodsInNamespaceMatchingSelector(selector, ns interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllPodsInNamespaceMatchingSelector", reflect.TypeOf((*MockClient)(nil).GetAllPodsInNamespaceMatchingSelector), selector, ns)
}

// GetAllResourcesFromSelector mocks base method.
func (m *MockClient) GetAllResourcesFromSelector(selector, ns string) ([]unstructured.Unstructured, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllResourcesFromSelector", selector, ns)
	ret0, _ := ret[0].([]unstructured.Unstructured)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllResourcesFromSelector indicates an expected call of GetAllResourcesFromSelector.
func (mr *MockClientMockRecorder) GetAllResourcesFromSelector(selector, ns interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllResourcesFromSelector", reflect.TypeOf((*MockClient)(nil).GetAllResourcesFromSelector), selector, ns)
}

// GetPodLogs mocks base method.
func (m *MockClient) GetPodLogs(podName, containerName string, followLog bool) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPodLogs", podName, containerName, followLog)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPodLogs indicates an expected call of GetPodLogs.
func (mr *MockClientMockRecorder) GetPodLogs(podName, containerName, followLog interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPodLogs", reflect.TypeOf((*MockClient)(nil).GetPodLogs), podName, containerName, followLog)
}

// GetPodsMatchingSelector mocks base method.
func (m *MockClient) GetPodsMatchingSelector(selector string) (*v1.PodList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPodsMatchingSelector", selector)
	ret0, _ := ret[0].(*v1.PodList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPodsMatchingSelector indicates an expected call of GetPodsMatchingSelector.
func (mr *MockClientMockRecorder) GetPodsMatchingSelector(selector interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPodsMatchingSelector", reflect.TypeOf((*MockClient)(nil).GetPodsMatchingSelector), selector)
}

// GetRunningPodFromSelector mocks base method.
func (m *MockClient) GetRunningPodFromSelector(selector string) (*v1.Pod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRunningPodFromSelector", selector)
	ret0, _ := ret[0].(*v1.Pod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRunningPodFromSelector indicates an expected call of GetRunningPodFromSelector.
func (mr *MockClientMockRecorder) GetRunningPodFromSelector(selector interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRunningPodFromSelector", reflect.TypeOf((*MockClient)(nil).GetRunningPodFromSelector), selector)
}

// KubeGenerate mocks base method.
func (m *MockClient) KubeGenerate(name string) (*v1.Pod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KubeGenerate", name)
	ret0, _ := ret[0].(*v1.Pod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// KubeGenerate indicates an expected call of KubeGenerate.
func (mr *MockClientMockRecorder) KubeGenerate(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KubeGenerate", reflect.TypeOf((*MockClient)(nil).KubeGenerate), name)
}

// ListAllComponents mocks base method.
func (m *MockClient) ListAllComponents() ([]api.ComponentAbstract, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAllComponents")
	ret0, _ := ret[0].([]api.ComponentAbstract)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAllComponents indicates an expected call of ListAllComponents.
func (mr *MockClientMockRecorder) ListAllComponents() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAllComponents", reflect.TypeOf((*MockClient)(nil).ListAllComponents))
}

// PlayKube mocks base method.
func (m *MockClient) PlayKube(pod *v1.Pod) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlayKube", pod)
	ret0, _ := ret[0].(error)
	return ret0
}

// PlayKube indicates an expected call of PlayKube.
func (mr *MockClientMockRecorder) PlayKube(pod interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlayKube", reflect.TypeOf((*MockClient)(nil).PlayKube), pod)
}

// PodLs mocks base method.
func (m *MockClient) PodLs() (map[string]bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PodLs")
	ret0, _ := ret[0].(map[string]bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PodLs indicates an expected call of PodLs.
func (mr *MockClientMockRecorder) PodLs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PodLs", reflect.TypeOf((*MockClient)(nil).PodLs))
}

// PodRm mocks base method.
func (m *MockClient) PodRm(podname string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PodRm", podname)
	ret0, _ := ret[0].(error)
	return ret0
}

// PodRm indicates an expected call of PodRm.
func (mr *MockClientMockRecorder) PodRm(podname interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PodRm", reflect.TypeOf((*MockClient)(nil).PodRm), podname)
}

// PodStop mocks base method.
func (m *MockClient) PodStop(podname string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PodStop", podname)
	ret0, _ := ret[0].(error)
	return ret0
}

// PodStop indicates an expected call of PodStop.
func (mr *MockClientMockRecorder) PodStop(podname interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PodStop", reflect.TypeOf((*MockClient)(nil).PodStop), podname)
}

// VolumeLs mocks base method.
func (m *MockClient) VolumeLs() (map[string]bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VolumeLs")
	ret0, _ := ret[0].(map[string]bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VolumeLs indicates an expected call of VolumeLs.
func (mr *MockClientMockRecorder) VolumeLs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VolumeLs", reflect.TypeOf((*MockClient)(nil).VolumeLs))
}

// VolumeRm mocks base method.
func (m *MockClient) VolumeRm(volumeName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VolumeRm", volumeName)
	ret0, _ := ret[0].(error)
	return ret0
}

// VolumeRm indicates an expected call of VolumeRm.
func (mr *MockClientMockRecorder) VolumeRm(volumeName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VolumeRm", reflect.TypeOf((*MockClient)(nil).VolumeRm), volumeName)
}
