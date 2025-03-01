package integration

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-developer/odo/tests/helper"
)

var _ = Describe("odo describe component command tests", func() {
	var commonVar helper.CommonVar
	var cmpName string

	// This is run before every Spec (It)
	var _ = BeforeEach(func() {
		commonVar = helper.CommonBeforeEach()
		helper.Chdir(commonVar.Context)
		cmpName = helper.RandString(6)
	})

	// This is run after every Spec (It)
	var _ = AfterEach(func() {
		helper.CommonAfterEach(commonVar)
	})

	for _, label := range []string{
		helper.LabelNoCluster, helper.LabelUnauth,
	} {
		label := label
		It("should fail, without cluster", Label(label), func() {
			By("running odo describe component -o json with namespace flag without name flag", func() {
				res := helper.Cmd("odo", "describe", "component", "--namespace", "default", "-o", "json").ShouldFail()
				stdout, stderr := res.Out(), res.Err()
				Expect(helper.IsJSON(stderr)).To(BeTrue())
				Expect(stdout).To(BeEmpty())
				helper.JsonPathContentContain(stderr, "message", "--namespace can be used only with --name")
			})

			By("running odo describe component -o json without name and without devfile in the current directory", func() {
				res := helper.Cmd("odo", "describe", "component", "-o", "json").ShouldFail()
				stdout, stderr := res.Out(), res.Err()
				Expect(helper.IsJSON(stderr)).To(BeTrue())
				Expect(stdout).To(BeEmpty())
				helper.JsonPathContentContain(stderr, "message", "The current directory does not represent an odo component")
			})

			By("running odo describe component with namespace flag without name flag", func() {
				res := helper.Cmd("odo", "describe", "component", "--namespace", "default").ShouldFail()
				stdout, stderr := res.Out(), res.Err()
				Expect(stdout).To(BeEmpty())
				Expect(stderr).To(ContainSubstring("--namespace can be used only with --name"))
			})

			By("running odo describe component without name and without devfile in the current directory", func() {
				res := helper.Cmd("odo", "describe", "component").ShouldFail()
				stdout, stderr := res.Out(), res.Err()
				Expect(stdout).To(BeEmpty())
				Expect(stderr).To(ContainSubstring("The current directory does not represent an odo component"))
			})

		})
	}

	It("should fail, with cluster", func() {
		By("running odo describe component -o json with an unknown name", func() {
			res := helper.Cmd("odo", "describe", "component", "--name", "unknown-name", "-o", "json").ShouldFail()
			stdout, stderr := res.Out(), res.Err()
			Expect(helper.IsJSON(stderr)).To(BeTrue())
			Expect(stdout).To(BeEmpty())
			helper.JsonPathContentContain(stderr, "message", "no component found with name \"unknown-name\" in the namespace \""+commonVar.Project+"\"")
		})

		By("running odo describe component with an unknown name", func() {
			res := helper.Cmd("odo", "describe", "component", "--name", "unknown-name").ShouldFail()
			stdout, stderr := res.Out(), res.Err()
			Expect(stdout).To(BeEmpty())
			Expect(stderr).To(ContainSubstring("no component found with name \"unknown-name\" in the namespace \"" + commonVar.Project + "\""))
		})
	})

	When("creating a component", func() {
		BeforeEach(func() {
			helper.Cmd("odo", "init", "--name", cmpName, "--devfile-path", helper.GetExamplePath("source", "devfiles", "nodejs", "devfile-registry.yaml")).ShouldPass()
		})

		checkDevfileJSONDescription := func(jsonContent string, devfileName string) {
			helper.JsonPathContentIs(jsonContent, "devfilePath", filepath.Join(commonVar.Context, devfileName))
			helper.JsonPathContentIs(jsonContent, "devfileData.devfile.metadata.name", cmpName)
			helper.JsonPathContentIs(jsonContent, "devfileData.supportedOdoFeatures.dev", "true")
			helper.JsonPathContentIs(jsonContent, "devfileData.supportedOdoFeatures.deploy", "false")
			helper.JsonPathContentIs(jsonContent, "devfileData.supportedOdoFeatures.debug", "true")
			helper.JsonPathContentIs(jsonContent, "managedBy", "odo")
		}

		checkDevfileDescription := func(content string, withUnknown bool) {
			Expect(content).To(ContainSubstring("Name: " + cmpName))
			Expect(content).To(ContainSubstring("Project Type: nodejs"))
			if withUnknown {
				for _, v := range []string{"Version", "Display Name", "Description", "Language"} {
					Expect(content).To(ContainSubstring(v + ": Unknown"))
				}
			}
		}

		for _, label := range []string{
			helper.LabelNoCluster, helper.LabelUnauth,
		} {
			label := label

			It("should describe the component in the current directory", Label(label), func() {
				By("running with json output", func() {
					res := helper.Cmd("odo", "describe", "component", "-o", "json").ShouldPass()
					stdout, stderr := res.Out(), res.Err()
					Expect(helper.IsJSON(stdout)).To(BeTrue())
					Expect(stderr).To(BeEmpty())
					checkDevfileJSONDescription(stdout, "devfile.yaml")
					helper.JsonPathContentIs(stdout, "runningIn", "")
					helper.JsonPathContentIs(stdout, "devForwardedPorts", "")
				})

				By("running with default output", func() {
					res := helper.Cmd("odo", "describe", "component").ShouldPass()
					stdout := res.Out()
					checkDevfileDescription(stdout, false)
					Expect(stdout).To(ContainSubstring("Running in: None"))
					Expect(stdout).ToNot(ContainSubstring("Forwarded ports"))
				})
			})

			When("renaming to hide devfile.yaml file", Label(label), func() {
				BeforeEach(func() {
					err := os.Rename("devfile.yaml", ".devfile.yaml")
					Expect(err).NotTo(HaveOccurred())
				})

				It("should describe the component in the current directory using the hidden devfile", func() {
					By("running with json output", func() {
						res := helper.Cmd("odo", "describe", "component", "-o", "json").ShouldPass()
						stdout, stderr := res.Out(), res.Err()
						Expect(helper.IsJSON(stdout)).To(BeTrue())
						Expect(stderr).To(BeEmpty())
						checkDevfileJSONDescription(stdout, ".devfile.yaml")
						helper.JsonPathContentIs(stdout, "runningIn", "")
						helper.JsonPathContentIs(stdout, "devForwardedPorts", "")
					})

					By("running with default output", func() {
						res := helper.Cmd("odo", "describe", "component").ShouldPass()
						stdout := res.Out()
						checkDevfileDescription(stdout, false)
						Expect(stdout).To(ContainSubstring("Running in: None"))
						Expect(stdout).ToNot(ContainSubstring("Forwarded ports"))
					})
				})
			})
		}

		It("should not describe the component from another directory", func() {
			By("running with json output", func() {
				err := os.Chdir("/")
				Expect(err).NotTo(HaveOccurred())
				res := helper.Cmd("odo", "describe", "component", "--name", cmpName, "-o", "json").ShouldFail()
				stdout, stderr := res.Out(), res.Err()
				Expect(helper.IsJSON(stderr)).To(BeTrue())
				Expect(stdout).To(BeEmpty())
				helper.JsonPathContentContain(stderr, "message", "no component found with name \""+cmpName+"\" in the namespace \""+commonVar.Project+"\"")
			})

			By("running with default output", func() {
				err := os.Chdir("/")
				Expect(err).NotTo(HaveOccurred())
				res := helper.Cmd("odo", "describe", "component", "--name", cmpName).ShouldFail()
				stdout, stderr := res.Out(), res.Err()
				Expect(stdout).To(BeEmpty())
				Expect(stderr).To(ContainSubstring("no component found with name \"" + cmpName + "\" in the namespace \"" + commonVar.Project + "\""))
			})
		})

		When("running odo dev", func() {
			var devSession helper.DevSession
			var ports map[string]string

			BeforeEach(func() {
				var err error
				devSession, _, _, ports, err = helper.StartDevMode(helper.DevSessionOpts{})
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				devSession.Kill()
				devSession.WaitEnd()
			})

			It("should describe the component in dev mode", func() {
				By("running with json output", func() {
					res := helper.Cmd("odo", "describe", "component", "-o", "json").ShouldPass()
					stdout, stderr := res.Out(), res.Err()
					Expect(helper.IsJSON(stdout)).To(BeTrue())
					Expect(stderr).To(BeEmpty())
					checkDevfileJSONDescription(stdout, "devfile.yaml")
					helper.JsonPathContentIs(stdout, "devForwardedPorts.#", "1")
					helper.JsonPathContentIs(stdout, "devForwardedPorts.0.containerName", "runtime")
					helper.JsonPathContentIs(stdout, "devForwardedPorts.0.localAddress", "127.0.0.1")
					helper.JsonPathContentIs(stdout, "devForwardedPorts.0.localPort", ports["3000"][len("127.0.0.1:"):])
					helper.JsonPathContentIs(stdout, "devForwardedPorts.0.containerPort", "3000")
				})

				By("running with default output", func() {
					res := helper.Cmd("odo", "describe", "component").ShouldPass()
					stdout := res.Out()
					checkDevfileDescription(stdout, false)
					Expect(stdout).To(ContainSubstring("Forwarded ports"))
					Expect(stdout).To(ContainSubstring("127.0.0.1:" + ports["3000"][len("127.0.0.1:"):] + " -> runtime:3000"))
				})
			})

			It("should describe the component from another directory", func() {
				By("running with json output", func() {
					err := os.Chdir("/")
					Expect(err).NotTo(HaveOccurred())
					res := helper.Cmd("odo", "describe", "component", "--name", cmpName, "-o", "json").ShouldPass()
					stdout, stderr := res.Out(), res.Err()
					Expect(helper.IsJSON(stdout)).To(BeTrue())
					Expect(stderr).To(BeEmpty())
					helper.JsonPathContentIs(stdout, "devfilePath", "")
					helper.JsonPathContentIs(stdout, "devfileData.devfile.metadata.name", cmpName)
					helper.JsonPathContentIs(stdout, "devfileData.devfile.metadata.projectType", "nodejs")
					for _, v := range []string{"version", "displayName", "description", "language"} {
						helper.JsonPathContentIs(stdout, "devfileData.devfile.metadata."+v, "Unknown")
					}
					helper.JsonPathContentIs(stdout, "devForwardedPorts", "")
					helper.JsonPathContentIs(stdout, "runningIn.dev", "true")
					helper.JsonPathContentIs(stdout, "runningIn.deploy", "false")
					helper.JsonPathContentIs(stdout, "devForwardedPorts", "")
				})

				By("running with default output", func() {
					err := os.Chdir("/")
					Expect(err).NotTo(HaveOccurred())
					res := helper.Cmd("odo", "describe", "component", "--name", cmpName).ShouldPass()
					stdout := res.Out()
					checkDevfileDescription(stdout, true)
					Expect(stdout).ToNot(ContainSubstring("Forwarded ports"))
					Expect(stdout).To(ContainSubstring("Running in: Dev"))
					Expect(stdout).To(ContainSubstring("Dev: Unknown"))
					Expect(stdout).To(ContainSubstring("Deploy: Unknown"))
					Expect(stdout).To(ContainSubstring("Debug: Unknown"))
				})
			})

		})

		for _, ctx := range []struct {
			title           string
			devfile         string
			matchOutput     []string
			matchJSONOutput map[string]string
		}{
			{
				title: "ingress/routes",
				devfile: func() string {
					if helper.IsKubernetesCluster() {
						return "devfile-deploy-ingress.yaml"
					}
					return "devfile-deploy-route.yaml"
				}(),
				matchOutput: func() []string {
					if helper.IsKubernetesCluster() {
						return []string{"Kubernetes Ingresses", "nodejs.example.com/", "nodejs.example.com/foo"}
					}
					return []string{"OpenShift Routes", "/foo"}
				}(),
				matchJSONOutput: func() map[string]string {
					if helper.IsKubernetesCluster() {
						return map[string]string{"ingresses.0.name": "my-nodejs-app", "ingresses.0.rules.0.host": "nodejs.example.com", "ingresses.0.rules.0.paths.0": "/", "ingresses.0.rules.0.paths.1": "/foo"}
					}
					return map[string]string{"routes.0.name": "my-nodejs-app", "routes.0.rules.0.paths.0": "/foo"}
				}(),
			},
			{
				title:           "ingress with defaultBackend",
				devfile:         "devfile-deploy-defaultBackend-ingress.yaml",
				matchOutput:     []string{"Kubernetes Ingresses", "*/*"},
				matchJSONOutput: map[string]string{"ingresses.0.name": "my-nodejs-app", "ingresses.0.rules.0.host": "*", "ingresses.0.rules.0.paths.0": "/*"},
			},
		} {
			ctx := ctx
			When("running odo deploy to create ingress/routes", func() {
				const (
					componentName = "nodejs-prj1-api-abhz" // hard-coded from the Devfiles
				)
				BeforeEach(func() {
					helper.CopyExampleDevFile(filepath.Join("source", "devfiles", "nodejs", ctx.devfile), path.Join(commonVar.Context, "devfile.yaml"))
					helper.Cmd("odo", "deploy").AddEnv("PODMAN_CMD=echo").ShouldPass()
				})
				It(fmt.Sprintf("should show the %s in odo describe component output", ctx.title), func() {
					By("checking the human readable output", func() {
						out := helper.Cmd("odo", "describe", "component").ShouldPass().Out()
						helper.MatchAllInOutput(out, ctx.matchOutput)
					})
					By("checking the machine readable output", func() {
						out := helper.Cmd("odo", "describe", "component", "-o", "json").ShouldPass().Out()
						for key, value := range ctx.matchJSONOutput {
							helper.JsonPathContentContain(out, key, value)
						}
					})
					By("checking the human readable output with component name", func() {
						out := helper.Cmd("odo", "describe", "component", "--name", componentName).ShouldPass().Out()
						helper.MatchAllInOutput(out, ctx.matchOutput)
					})
					By("checking the machine readable output with component name", func() {
						out := helper.Cmd("odo", "describe", "component", "--name", componentName, "-o", "json").ShouldPass().Out()
						for key, value := range ctx.matchJSONOutput {
							helper.JsonPathContentContain(out, key, value)
						}
					})
				})
			})
		}
	})
})
