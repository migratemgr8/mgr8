package applications

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type SectionProperty struct {
	key   string
	value string
}

type ConfigurationService struct{}

type Configuration struct {
	Username string
	Hostname string
}

const (
	ConfigFileName   string = ".mgr8config"
	RunWithDockerEnv string = "RUN_WITH_DOCKER"
	UsernameEnv      string = "MGR8_USERNAME"
	UserSection      string = "user"
	UsernameKey      string = "username"
	HostnameKey      string = "hostname"
)

func NewConfigurationService() *ConfigurationService {
	return &ConfigurationService{}
}

func CreateConfiguration() (Configuration, error) {
	configuration := Configuration{}

	configFilePath, err := GetConfigFilePath()
	if err != nil {
		return configuration, err
	}

	configurationFile, err := os.Create(configFilePath)
	if err != nil {
		return configuration, err
	}

	fmt.Println("Configuration file not found. Configure:")

	// configuring user details (username and hostname)
	configuration.Hostname, err = os.Hostname()
	if err != nil {
		return configuration, err
	}

	configuration.Username = configuration.Hostname

	fmt.Println("Your default username is " + configuration.Username + ". It will be displayed on the logs when you execute a migration.")
	var answer byte

	for !IsValidAnswer(answer) {
		fmt.Println("Do you want to change it? (y/n)")
		fmt.Scanf("%c ", &answer)
	}

	if IsYesAnswer(answer) {
		isValidUsername := false
		scanner := bufio.NewScanner(os.Stdin)

		username := ""

		for !isValidUsername {
			if err != nil {
				fmt.Println("This username is not valid: " + err.Error())
			}
			fmt.Println("Please enter your username:")
			scanner.Scan()
			if err := scanner.Err(); err != nil {
				return configuration, err
			}
			username = scanner.Text()
			isValidUsername, err = IsValidUsername(username)
		}

		configuration.Username = username
	}

	if err = InsertUserDetails(configuration.Username, configuration.Hostname, configurationFile); err != nil {
		return configuration, err
	}

	if err = configurationFile.Close(); err != nil {
		return configuration, err
	}

	fmt.Println("MGR8 Configured successfuly!")

	return configuration, err
}

func (c *ConfigurationService) GetConfigurations() (Configuration, error) {
	configuration := Configuration{}

	// if running with docker use env configuration
	if os.Getenv(RunWithDockerEnv) != "" {
		usernameEnv := os.Getenv(UsernameEnv)
		if usernameEnv == "" {
			return configuration, errors.New("no username was found, set it on env " + UsernameEnv)
		}

		configuration.Username = usernameEnv
		configuration.Hostname = usernameEnv

		return configuration, nil
	}

	configFilePath, err := GetConfigFilePath()
	if err != nil {
		return configuration, err
	}

	configurationFile, err := os.Open(configFilePath)

	// if configuration file doesnt exist
	if errors.Is(err, os.ErrNotExist) {
		configuration, err = CreateConfiguration()
		return configuration, err
	}

	scanner := bufio.NewScanner(configurationFile)

	var section string
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return configuration, err
		}

		if GetSection(scanner.Text()) != "" {
			section = GetSection(scanner.Text())
			continue
		}

		property := GetSectionProperty(scanner.Text())

		if section == UserSection {
			if property.key == UsernameKey {
				configuration.Username = property.value
			} else if property.key == HostnameKey {
				configuration.Hostname = property.value
			}
		}
	}

	if err = configurationFile.Close(); err != nil {
		return configuration, err
	}

	return configuration, err
}

func GetSection(line string) string {
	re := regexp.MustCompile(`\[(.*?)\]`)
	matches := re.FindAllStringSubmatch(line, 1)

	if matches == nil {
		return ""
	}

	section := matches[0][1]

	return section
}

func GetSectionProperty(line string) SectionProperty {
	re := regexp.MustCompile(`(\S*)\s*=\s*(\S*)`)
	matches := re.FindAllStringSubmatch(line, 1)

	if matches == nil {
		return SectionProperty{}
	}

	property := SectionProperty{
		key:   matches[0][1],
		value: matches[0][2],
	}

	return property
}

func GetConfigFilePath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFilePath := filepath.Join(userHomeDir, ConfigFileName)

	return configFilePath, err
}

func IsYesAnswer(answer byte) bool {
	return answer == 'y' || answer == 'Y'
}

func IsValidAnswer(answer byte) bool {
	return answer == 'y' || answer == 'Y' || answer == 'n' || answer == 'N'
}

func IsValidUsername(username string) (bool, error) {
	if username == "" {
		return false, errors.New("username cannot be empty")
	}
	if strings.ContainsRune(username, ' ') {
		return false, errors.New("username cannot contain spaces")
	}

	return true, nil
}

func InsertUserDetails(username string, hostname string, config *os.File) error {
	_, err := config.WriteString("[" + UserSection + "]\n\t" + UsernameKey + " = " + username + "\n\t" + HostnameKey + " = " + hostname + "\n")
	return err
}
