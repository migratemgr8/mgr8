package applications

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

type UserDetails struct {
	username string
	hostname string
}

const (
	ConfigurationFilename string = ".mgr8config"
	RunWithDockerEnv      string = "RUN_WITH_DOCKER"
	UsernameEnv           string = "MGR8_USERNAME"
	userSection           string = "user"
	usernameKey           string = "username"
	hostnameKey           string = "hostname"
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

	fmt.Println("Configuration file not found. Configure the tool:")

	// configure user details
	userDetails, err := GetUserDetails()
	if err != nil {
		return configuration, err
	}

	configuration.Username = userDetails.username
	configuration.Hostname = userDetails.hostname

	if err = InsertUserDetails(userDetails.username, userDetails.hostname, configurationFile); err != nil {
		return configuration, err
	}

	if err = configurationFile.Close(); err != nil {
		return configuration, err
	}

	fmt.Println("MGR8 Configured successfuly!")

	return configuration, err
}

func GetUserDetails() (UserDetails, error) {
	userDetails := UserDetails{}

	// configuring user details (username and hostname)
	hostname, err := os.Hostname()
	if err != nil {
		return userDetails, err
	}

	userDetails.hostname = hostname
	userDetails.username = userDetails.hostname

	fmt.Println("Your default username is " + userDetails.username + ". It will be displayed on the logs when you execute a migration.")
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
				return userDetails, err
			}
			username = scanner.Text()
			isValidUsername, err = IsValidUsername(username)
		}

		userDetails.username = username
	}

	return userDetails, err
}

func (c *ConfigurationService) GetConfigurations() (Configuration, error) {
	configuration := Configuration{}

	// if running with docker use env configuration
	if os.Getenv(RunWithDockerEnv) == "true" {
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

	// if configuration file doesn't exist
	if errors.Is(err, os.ErrNotExist) {
		configuration, err = CreateConfiguration()
		return configuration, err
	}

	configuration, err = ParseConfiguration(configurationFile)

	if err := configurationFile.Close(); err != nil {
		return configuration, err
	}

	return configuration, err
}

func ParseConfiguration(configurationFile io.Reader) (Configuration, error) {
	configuration := Configuration{}
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

		if section == userSection {
			if property.key == usernameKey {
				configuration.Username = property.value
			} else if property.key == hostnameKey {
				configuration.Hostname = property.value
			}
		}
	}

	return configuration, nil
}

func GetConfigFilePath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(userHomeDir, ConfigurationFilename)

	return path, err
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
	_, err := config.WriteString("[" + userSection + "]\n\t" + usernameKey + " = " + username + "\n\t" + hostnameKey + " = " + hostname + "\n")
	return err
}
