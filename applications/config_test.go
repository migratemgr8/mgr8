package applications

import (
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration", func() {
	var (
		subject *ConfigurationService = NewConfigurationService()
	)

	Context("Configuration Service", func() {

		When("GetConfigurations is called", func() {

			When("RunWithDockerEnv is set", func() {
				BeforeEach(func() {
					os.Setenv(RunWithDockerEnv, "true")
				})

				It("Returns empty username and error", func() {
					configuration, err := subject.GetConfigurations()

					Expect(err).ToNot(BeNil())
					Expect(configuration).To(Equal(Configuration{}))
				})

				When("UsernameEnv is set", func() {
					BeforeEach(func() {
						os.Setenv(UsernameEnv, "docker_username")
					})

					It("Returns UsernameEnv as username", func() {
						configuration, err := subject.GetConfigurations()

						Expect(err).To(BeNil())
						Expect(configuration).To(Equal(Configuration{
							Username: "docker_username",
							Hostname: "docker_username",
						}))
					})
				})

				AfterSuite(func() {
					// reset envs for next tests
					os.Setenv(RunWithDockerEnv, "")
					os.Setenv(UsernameEnv, "")
				})
			})

			FWhen("RunWithDockerEnv is not set", func() {
				BeforeEach(func() {
					userHomeDir, err := os.UserHomeDir()
					configFilePath := filepath.Join(userHomeDir, ConfigurationFilename)
					_, err = os.Stat(configFilePath)
					if errors.Is(err, os.ErrNotExist) {
						os.Remove(configFilePath)
					}
				})

				// FWhen("Configuration file is not found", func() {
				// 	configuration, err := subject.GetConfigurations()

				// 	FIt("", func() {

				// 	})
				// })

				// FWhen("Configuration file is found", func() {
				// 	FIt("", func() {

				// 	})
				// })
			})
		})

	})
})
