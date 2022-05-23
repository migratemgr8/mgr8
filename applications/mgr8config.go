package applications

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type userNameService struct{}

func NewUserNameService() *userNameService {
	return &userNameService{}
}

func (a *userNameService) GetUserName() (string, error) {

	// trying to read .mgr8config file if exists
	user_home_directory, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	MGR8CONFIG_FILEPATH := user_home_directory + "/.mgr8config"
	mgr8config, err := os.Open(MGR8CONFIG_FILEPATH)

	// close the file with defer
	defer mgr8config.Close()
	var username string

	if err == nil {
		// get username from .mgr8config
		scanner := bufio.NewScanner(mgr8config)
		// reading file line by line
		for scanner.Scan() {
			key, value := GetKeyValueFromMgr8Config(scanner.Text())
			if key == "name" {
				username = value
			}
		}

		if err := scanner.Err(); err != nil {
			return "", err
		}

		return username, nil
	} else {
		// create .mgr8config file
		mgr8config, err := os.Create(MGR8CONFIG_FILEPATH)
		if err != nil {
			return "", err
		}

		// close the file with defer
		defer mgr8config.Close()

		hostname, err := os.Hostname()
		if err != nil {
			return "", err
		}

		fmt.Println("File .mgr8config not found. Please enter your username (without whitespaces):")
		fmt.Scanf("%s", &username)

		// write a string on file
		mgr8config.WriteString("[user]\n\tname = " + username + "\n\tdevice_name = " + hostname + "\n")

		return username, nil
	}
}

func GetKeyValueFromMgr8Config(mgr8configLine string) (string, string) {
	mgr8configLine = strings.Trim(mgr8configLine, "\t")
	key_and_value := strings.Split(mgr8configLine, " = ")
	key, value := key_and_value[0], key_and_value[len(key_and_value)-1]
	return key, value
}
