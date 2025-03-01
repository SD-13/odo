package integration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/library/pkg/devfile/parser"
	"github.com/devfile/library/pkg/devfile/parser/data/v2/common"
	"k8s.io/utils/pointer"

	"github.com/redhat-developer/odo/pkg/odo/cli/messages"
	"github.com/redhat-developer/odo/pkg/util"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	odolog "github.com/redhat-developer/odo/pkg/log"
	"github.com/redhat-developer/odo/pkg/version"
	"github.com/redhat-developer/odo/tests/helper"
)

var _ = Describe("odo init interactive command tests", func() {
	for _, label := range []string{
		helper.LabelNoCluster, helper.LabelUnauth,
	} {
		label := label
		var _ = Context("label "+label, Label(label), func() {
			var commonVar helper.CommonVar

			// This is run before every Spec (It)
			var _ = BeforeEach(func() {
				commonVar = helper.CommonBeforeEach()
				helper.Chdir(commonVar.Context)

				// We make EXPLICITLY sure that we are outputting with NO COLOR
				// this is because in some cases we are comparing the output with a colorized one
				os.Setenv("NO_COLOR", "true")
			})

			// Clean up after the test
			// This is run after every Spec (It)
			var _ = AfterEach(func() {
				helper.CommonAfterEach(commonVar)
			})

			It("should not fail when using -v flag", func() {
				command := []string{"odo", "init", "-v", "4"}
				output, err := helper.RunInteractive(command, nil, func(ctx helper.InteractiveContext) {

					By("showing the interactive mode notice message", func() {
						helper.ExpectString(ctx, messages.InteractiveModeEnabled)
					})

					helper.ExpectString(ctx, "Select language")
					helper.SendLine(ctx, "Go")

					helper.ExpectString(ctx, "Select project type")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Select version")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Select container for which you want to change configuration?")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Which starter project do you want to use")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Enter component name")
					helper.SendLine(ctx, "my-go-app")

					helper.ExpectString(ctx, "Your new component 'my-go-app' is ready in the current directory")

				})
				Expect(err).To(BeNil())
				Expect(output).To(ContainSubstring("Your new component 'my-go-app' is ready in the current directory"))
				Expect(helper.ListFilesInDir(commonVar.Context)).To(ContainElements("devfile.yaml"))
			})

			It("should ask to re-enter the component name when an invalid value is passed", func() {
				command := []string{"odo", "init"}
				_, err := helper.RunInteractive(command, nil, func(ctx helper.InteractiveContext) {

					helper.ExpectString(ctx, "Select language")
					helper.SendLine(ctx, "Go")

					helper.ExpectString(ctx, "Select project type")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Select version")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Select container for which you want to change configuration?")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Which starter project do you want to use")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Enter component name")
					helper.SendLine(ctx, "myapp-<script>alert('Injected!');</script>")

					helper.ExpectString(ctx, "name \"myapp-<script>alert('Injected!');</script>\" is not valid, name should conform the following requirements")

					helper.ExpectString(ctx, "Enter component name")
					helper.SendLine(ctx, "my-go-app")

					helper.ExpectString(ctx, "Your new component 'my-go-app' is ready in the current directory")
				})
				Expect(err).To(BeNil())
				Expect(helper.ListFilesInDir(commonVar.Context)).To(ContainElements("devfile.yaml"))
			})

			It("should print automation command with proper values", func() {

				// This test fails on Windows because of terminal emulator behaviour
				if os.Getenv("SKIP_WELCOMING_MESSAGES") == "true" {
					Skip("This is a Unix specific scenario, skipping")
				}

				command := []string{"odo", "init"}
				starter := "go-starter"
				componentName := "my-go-app"
				devfileName := "go"
				devfileVersion := "1.0.2"

				output, err := helper.RunInteractive(command, nil, func(ctx helper.InteractiveContext) {

					helper.ExpectString(ctx, "Select language")
					helper.SendLine(ctx, "Go")

					helper.ExpectString(ctx, "Select project type")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Select version")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Select container for which you want to change configuration?")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Which starter project do you want to use")
					helper.SendLine(ctx, starter)

					helper.ExpectString(ctx, "Enter component name")
					helper.SendLine(ctx, componentName)

					helper.ExpectString(ctx, "Your new component 'my-go-app' is ready in the current directory")

				})

				Expect(err).To(BeNil())
				Expect(output).To(ContainSubstring("odo init --name %s --devfile %s --devfile-registry DefaultDevfileRegistry --devfile-version %s --starter %s", componentName, devfileName, devfileVersion, starter))
				Expect(output).To(ContainSubstring("Your new component 'my-go-app' is ready in the current directory"))
				Expect(helper.ListFilesInDir(commonVar.Context)).To(ContainElements("devfile.yaml"))
			})

			It("should download correct devfile", func() {

				command := []string{"odo", "init"}
				output, err := helper.RunInteractive(command, nil, func(ctx helper.InteractiveContext) {

					By("showing the interactive mode notice message", func() {
						helper.ExpectString(ctx, messages.InteractiveModeEnabled)
					})

					helper.ExpectString(ctx, "Select language")
					helper.SendLine(ctx, "Go")

					helper.ExpectString(ctx, "Select project type")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Select version")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Select container for which you want to change configuration?")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Which starter project do you want to use")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Enter component name")
					helper.SendLine(ctx, "my-go-app")

					helper.ExpectString(ctx, "Your new component 'my-go-app' is ready in the current directory")

				})

				Expect(err).To(BeNil())
				Expect(output).To(ContainSubstring("Your new component 'my-go-app' is ready in the current directory"))
				Expect(helper.ListFilesInDir(commonVar.Context)).To(ContainElements("devfile.yaml"))
			})

			It("should download correct devfile-starter", func() {

				command := []string{"odo", "init"}
				output, err := helper.RunInteractive(command, nil, func(ctx helper.InteractiveContext) {

					By("showing the interactive mode notice message", func() {
						helper.ExpectString(ctx, messages.InteractiveModeEnabled)
					})

					helper.ExpectString(ctx, "Select language")
					helper.SendLine(ctx, "java")

					helper.ExpectString(ctx, "Select project type")
					helper.SendLine(ctx, "Vert.x Java")

					helper.ExpectString(ctx, "Select container for which you want to change configuration?")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Which starter project do you want to use")
					helper.SendLine(ctx, "vertx-cache-example-redhat")

					helper.ExpectString(ctx, "Enter component name")
					helper.SendLine(ctx, "my-app")

					helper.ExpectString(ctx, "Your new component 'my-app' is ready in the current directory")

				})

				Expect(err).To(BeNil())
				Expect(output).To(ContainSubstring("Downloading starter project \"vertx-cache-example-redhat\""))
			})

			Describe("displaying welcoming messages", func() {

				// testFunc is a function that returns a `Tester` function (intended to be used via `helper.RunInteractive`),
				// which first expects all messages in `welcomingMsgs` to be read from the console,
				// then runs an `additionalTester` and finally expects the asking of a component name
				// (based on the `language` specified)
				testFunc := func(language string, welcomingMsgs []string, additionalTester helper.Tester) helper.Tester {
					return func(ctx helper.InteractiveContext) {
						for _, msg := range welcomingMsgs {
							helper.ExpectString(ctx, msg)
						}

						if additionalTester != nil {
							additionalTester(ctx)
						}

						helper.ExpectString(ctx, "Enter component name")
						helper.SendLine(ctx, fmt.Sprintf("my-%s-app", language))

						helper.ExpectString(ctx,
							fmt.Sprintf("Your new component 'my-%s-app' is ready in the current directory", language))
					}
				}

				assertBehavior := func(language string, output string, err error, msgs []string, additionalAsserter func()) {
					Expect(err).To(BeNil())

					lines, err := helper.ExtractLines(output)
					if err != nil {
						log.Fatal(err)
					}
					Expect(len(lines)).To(BeNumerically(">", len(msgs)))
					Expect(lines[0:len(msgs)]).To(Equal(msgs))
					Expect(lines).To(
						ContainElement(fmt.Sprintf("Your new component 'my-%s-app' is ready in the current directory", strings.ToLower(language))))

					Expect(helper.ListFilesInDir(commonVar.Context)).To(ContainElements("devfile.yaml"))

					if additionalAsserter != nil {
						additionalAsserter()
					}
				}

				testRunner := func(language string, welcomingMsgs []string, tester helper.Tester) (string, error) {
					command := []string{"odo", "init"}
					return helper.RunInteractive(command,
						// Setting verbosity level to 0, because we would be asserting the welcoming message is the first
						// message displayed to the end user. So we do not want any potential debug lines to be printed first.
						// Using envvars here (and not via the -v flag), because of https://github.com/redhat-developer/odo/issues/5513
						[]string{"ODO_LOG_LEVEL=0"},
						testFunc(strings.ToLower(language), welcomingMsgs, tester))
				}

				When("directory is empty", func() {

					BeforeEach(func() {
						Expect(helper.ListFilesInDir(commonVar.Context)).To(HaveLen(0))
					})

					It("should display appropriate welcoming messages", func() {

						if os.Getenv("SKIP_WELCOMING_MESSAGES") == "true" {
							Skip("This is a Unix specific scenario, skipping")
						}

						language := "java"

						// The first output is welcoming message / paragraph / banner output
						welcomingMsgs := strings.Split(odolog.Stitle(messages.InitializingNewComponent, messages.NoSourceCodeDetected, "odo version: "+version.VERSION), "\n")

						output, err := testRunner(language, welcomingMsgs, func(ctx helper.InteractiveContext) {
							helper.ExpectString(ctx, "Select language")
							helper.SendLine(ctx, language)

							helper.ExpectString(ctx, "Select project type")
							helper.SendLine(ctx, "")

							helper.ExpectString(ctx, "Select container for which you want to change configuration?")
							helper.SendLine(ctx, "")

							helper.ExpectString(ctx, "Which starter project do you want to use")
							helper.SendLine(ctx, "")
						})

						assertBehavior(strings.ToLower(language), output, err, welcomingMsgs, nil)
					})
				})

				When("directory is not empty", func() {

					BeforeEach(func() {
						helper.CopyExample(filepath.Join("source", "python"), commonVar.Context)
						Expect(helper.ListFilesInDir(commonVar.Context)).To(
							SatisfyAll(
								HaveLen(2),
								ContainElements("requirements.txt", "wsgi.py")))
					})

					It("should display appropriate welcoming messages", func() {

						if os.Getenv("SKIP_WELCOMING_MESSAGES") == "true" {
							Skip("This is a Unix specific scenario, skipping")
						}

						language := "Python"
						projectType := "Python"
						versionedDevfileName := "python:2.1.0"
						welcomingMsgs := strings.Split(odolog.Stitle(messages.InitializingNewComponent, messages.SourceCodeDetected, "odo version: "+version.VERSION), "\n")

						output, err := testRunner(language, welcomingMsgs, func(ctx helper.InteractiveContext) {
							helper.ExpectString(ctx, "Based on the files in the current directory odo detected")

							helper.ExpectString(ctx, fmt.Sprintf("Language: %s", language))

							helper.ExpectString(ctx, fmt.Sprintf("Project type: %s", projectType))

							helper.ExpectString(ctx,
								fmt.Sprintf("The devfile %q from the registry \"DefaultDevfileRegistry\" will be downloaded.", versionedDevfileName))

							helper.ExpectString(ctx, "Is this correct")
							helper.SendLine(ctx, "")

							helper.ExpectString(ctx, "Select container for which you want to change configuration")
							helper.SendLine(ctx, "")
						})

						assertBehavior(strings.ToLower(language), output, err, welcomingMsgs, func() {
							// Make sure the original source code files are still present
							Expect(helper.ListFilesInDir(commonVar.Context)).To(
								SatisfyAll(
									HaveLen(4),
									ContainElements("devfile.yaml", "requirements.txt", "wsgi.py", util.DotOdoDirectory)))
						})
					})
				})

				When("alizer detection of javascript name", func() {

					BeforeEach(func() {
						helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
						Expect(helper.ListFilesInDir(commonVar.Context)).To(
							SatisfyAll(
								HaveLen(3),
								ContainElements("Dockerfile", "package.json", "server.js")))
					})

					It("should display node-echo name", func() {
						language := "JavaScript"
						projectType := "Node.js"
						projectName := "node-echo"
						versionedDevfileName := "nodejs:2.1.1"

						output, err := helper.RunInteractive([]string{"odo", "init"}, nil, func(ctx helper.InteractiveContext) {
							helper.ExpectString(ctx, "Based on the files in the current directory odo detected")

							helper.ExpectString(ctx, fmt.Sprintf("Language: %s", language))

							helper.ExpectString(ctx, fmt.Sprintf("Project type: %s", projectType))

							helper.ExpectString(ctx,
								fmt.Sprintf("The devfile %q from the registry \"DefaultDevfileRegistry\" will be downloaded.", versionedDevfileName))

							helper.ExpectString(ctx, "Is this correct")
							helper.SendLine(ctx, "")

							helper.ExpectString(ctx, "Select container for which you want to change configuration")
							helper.SendLine(ctx, "")

							helper.ExpectString(ctx, "Enter component name")
							helper.SendLine(ctx, "")

							helper.ExpectString(ctx, fmt.Sprintf("Your new component '%s' is ready in the current directory", projectName))
						})
						Expect(err).To(BeNil())

						lines, err := helper.ExtractLines(output)
						Expect(err).To(BeNil())
						Expect(len(lines)).To(BeNumerically(">", 2))
						Expect(lines[len(lines)-1]).To(Equal(fmt.Sprintf("Your new component '%s' is ready in the current directory", projectName)))

					})
					It("should ask to re-enter the component name if invalid value is passed by the user", func() {
						language := "JavaScript"
						projectType := "Node.js"
						projectName := "node-echo"
						versionedDevfileName := "nodejs:2.1.1"

						_, err := helper.RunInteractive([]string{"odo", "init"}, nil, func(ctx helper.InteractiveContext) {
							helper.ExpectString(ctx, "Based on the files in the current directory odo detected")

							helper.ExpectString(ctx, fmt.Sprintf("Language: %s", language))

							helper.ExpectString(ctx, fmt.Sprintf("Project type: %s", projectType))

							helper.ExpectString(ctx,
								fmt.Sprintf("The devfile %q from the registry \"DefaultDevfileRegistry\" will be downloaded.", versionedDevfileName))

							helper.ExpectString(ctx, "Is this correct")
							helper.SendLine(ctx, "")

							helper.ExpectString(ctx, "Select container for which you want to change configuration")
							helper.SendLine(ctx, "")

							helper.ExpectString(ctx, "Enter component name")
							helper.SendLine(ctx, "myapp-<script>alert('Injected!');</script>")

							helper.ExpectString(ctx, "name \"myapp-<script>alert('Injected!');</script>\" is not valid, name should conform the following requirements")

							helper.ExpectString(ctx, "Enter component name")
							helper.SendLine(ctx, "")

							helper.ExpectString(ctx, fmt.Sprintf("Your new component '%s' is ready in the current directory", projectName))
						})
						Expect(err).To(BeNil())
						Expect(helper.ListFilesInDir(commonVar.Context)).To(ContainElements("devfile.yaml"))
					})
				})
			})

			It("should start downloading starter project only after all interactive questions have been asked", func() {

				output, err := helper.RunInteractive([]string{"odo", "init"}, nil, func(ctx helper.InteractiveContext) {

					helper.ExpectString(ctx, "Select language")
					helper.SendLine(ctx, ".NET")

					helper.ExpectString(ctx, "Select project type")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Select container for which you want to change configuration?")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Which starter project do you want to use")
					helper.SendLine(ctx, "")

					helper.ExpectString(ctx, "Enter component name")
					helper.SendLine(ctx, "my-dotnet-app")

					helper.ExpectString(ctx, "Your new component 'my-dotnet-app' is ready in the current directory")
				})

				Expect(err).To(BeNil())

				lines, err := helper.ExtractLines(output)
				Expect(err).To(BeNil())
				Expect(len(lines)).To(BeNumerically(">", 2))
				Expect(lines[len(lines)-1]).To(Equal("Your new component 'my-dotnet-app' is ready in the current directory"))

				componentNameQuestionIdx, ok := helper.FindFirstElementIndexMatchingRegExp(lines, ".*Enter component name:.*")
				Expect(ok).To(BeTrue())
				starterProjectDownloadActionIdx, found := helper.FindFirstElementIndexMatchingRegExp(lines,
					".*Downloading starter project \"([^\\s]+)\" \\[.*")
				Expect(found).To(BeTrue())
				Expect(starterProjectDownloadActionIdx).To(SatisfyAll(
					Not(BeZero()),
					// #5495: component name question should be displayed before starter project is actually downloaded
					BeNumerically(">", componentNameQuestionIdx),
				), "Action 'Downloading starter project' should have been displayed after the last interactive question ('Enter component name')")

				Expect(helper.ListFilesInDir(commonVar.Context)).To(ContainElements("devfile.yaml"))
			})

			Context("Automatic port detection via Alizer", func() {

				When("starting with an existing project", func() {
					const appPort = 34567

					BeforeEach(func() {
						helper.CopyExample(filepath.Join("source", "nodejs"), commonVar.Context)
						helper.ReplaceString(filepath.Join(commonVar.Context, "Dockerfile"), "EXPOSE 8080", fmt.Sprintf("EXPOSE %d", appPort))
					})

					It("should display ports detected", func() {
						_, err := helper.RunInteractive([]string{"odo", "init"}, nil, func(ctx helper.InteractiveContext) {
							helper.ExpectString(ctx, fmt.Sprintf("Application ports: %d", appPort))

							helper.SendLine(ctx, "Is this correct")
							helper.SendLine(ctx, "")

							helper.ExpectString(ctx, "Select container for which you want to change configuration")
							helper.SendLine(ctx, "")

							helper.ExpectString(ctx, "Enter component name")
							helper.SendLine(ctx, "my-nodejs-app-with-port-detected")

							helper.ExpectString(ctx, "Your new component 'my-nodejs-app-with-port-detected' is ready in the current directory")
						})
						Expect(err).ShouldNot(HaveOccurred())

						// Now make sure the Devfile contains a single container component with the right endpoint
						d, err := parser.ParseDevfile(parser.ParserArgs{Path: filepath.Join(commonVar.Context, "devfile.yaml"), FlattenedDevfile: pointer.BoolPtr(false)})
						Expect(err).ShouldNot(HaveOccurred())

						containerComponents, err := d.Data.GetDevfileContainerComponents(common.DevfileOptions{})
						Expect(err).ShouldNot(HaveOccurred())

						allPortsExtracter := func(comps []v1alpha2.Component) []int {
							var ports []int
							for _, c := range comps {
								for _, ep := range c.Container.Endpoints {
									ports = append(ports, ep.TargetPort)
								}
							}
							return ports
						}
						Expect(containerComponents).Should(WithTransform(allPortsExtracter, ContainElements(appPort)))
					})
				})
			})

			When("Alizer cannot determine a Devfile based on the current source code", func() {
				BeforeEach(func() {
					helper.CreateSimpleFile(commonVar.Context, "some-file-", ".ext")
				})

				It("should not fail but fallback to the interactive mode", func() {
					output, err := helper.RunInteractive([]string{"odo", "init"}, nil, func(ctx helper.InteractiveContext) {
						helper.ExpectString(ctx, "Could not determine a Devfile based on the files in the current directory")

						helper.ExpectString(ctx, "Select language")
						helper.SendLine(ctx, "PHP")

						helper.ExpectString(ctx, "Select project type")
						helper.SendLine(ctx, "")

						helper.ExpectString(ctx, "Select container for which you want to change configuration?")
						helper.SendLine(ctx, "")

						helper.ExpectString(ctx, "Enter component name")
						helper.SendLine(ctx, "my-php-app")

						helper.ExpectString(ctx, "Your new component 'my-php-app' is ready in the current directory")
					})

					Expect(err).ShouldNot(HaveOccurred())
					Expect(output).ShouldNot(ContainSubstring("Which starter project do you want to use"))
					Expect(helper.ListFilesInDir(commonVar.Context)).To(ContainElement("devfile.yaml"))
				})
			})
		})
	}
})
