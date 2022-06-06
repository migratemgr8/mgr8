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

type ConfigService struct {
	UserNameService UserNameService
}

type UserNameService struct{}

const (
	ConfigFileName string = ".mgr8config"
	UserNameEnv    string = "MGR8_USERNAME"
	UserSection    string = "user"
	UserNameKey    string = "username"
	HostNameKey    string = "hostname"
)

func NewUserNameService() *UserNameService {
	return &UserNameService{}
}

func (a *UserNameService) GetUserName() (string, error) {
	var username string

	// if env UserNameEnv was set, use it for username (priority, for use with docker)
	userNameEnv := os.Getenv(UserNameEnv)
	if userNameEnv != "" {
		return userNameEnv, nil
	}

	configFilePath, err := GetConfigFilePath()
	if err != nil {
		return "", err
	}

	config, err := os.Open(configFilePath)

	// if config not exists create it, configure it and return username
	if errors.Is(err, os.ErrNotExist) {
		config, err = os.Create(configFilePath)
		if err != nil {
			return "", err
		}

		fmt.Println("Configuration file not found. Configure:")

		hostname, err := os.Hostname()
		if err != nil {
			return "", err
		}

		username = hostname

		fmt.Println("Your default username is " + username + ". It will be displayed on the logs when you execute a migration.")
		var answer byte

		for !IsValidAnswer(answer) {
			fmt.Println("Do you want to change it? (y/n)")
			fmt.Scanf("%c ", &answer)
		}

		if IsYesAnswer(answer) {
			isValidUserName := false
			scanner := bufio.NewScanner(os.Stdin)

			for !isValidUserName {
				if err != nil {
					fmt.Println("This username is not valid: " + err.Error())
				}
				fmt.Println("Please enter your username:")
				scanner.Scan()
				if err := scanner.Err(); err != nil {
					return "", err
				}
				username = scanner.Text()
				isValidUserName, err = IsValidUserName(username)
			}
		}

		if err = InsertUserDetails(username, hostname, config); err != nil {
			return "", err
		}

		if err = config.Close(); err != nil {
			return "", err
		}

		fmt.Println("Username set")

		return username, err
	}

	// if config exists get username
	scanner := bufio.NewScanner(config)

	var section string
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}

		if GetSection(scanner.Text()) != "" {
			section = GetSection(scanner.Text())
			continue
		}

		property := GetSectionProperty(scanner.Text())

		if section == UserSection && property.key == UserNameKey {
			username = property.value
			break
		}
	}

	if err = config.Close(); err != nil {
		return "", err
	}

	return username, err
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
	user_home_dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFilePath := filepath.Join(user_home_dir, ConfigFileName)

	return configFilePath, err
}

func IsYesAnswer(answer byte) bool {
	return answer == 'y' || answer == 'Y'
}

func IsValidAnswer(answer byte) bool {
	return answer == 'y' || answer == 'Y' || answer == 'n' || answer == 'N'
}

func IsValidUserName(username string) (bool, error) {
	if username == "" {
		return false, errors.New("username cannot be empty")
	}
	if strings.ContainsRune(username, ' ') {
		return false, errors.New("username cannot contain spaces")
	}

	return true, nil
}

func InsertUserDetails(username string, hostname string, config *os.File) error {
	_, err := config.WriteString("[" + UserSection + "]\n\t" + UserNameKey + " = " + username + "\n\t" + HostNameKey + " = " + hostname + "\n")
	return err
}
