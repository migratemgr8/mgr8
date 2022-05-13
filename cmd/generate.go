package cmd

import (
	"fmt"
	"os"

	"github.com/kenji-yamane/mgr8/drivers"
)

type generate struct{}

func (g *generate) execute(args []string, databaseURL string, driver drivers.Driver) error {
	fileName := args[0]
	content, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("could not read from file: %s", err)
	}

	_, err = driver.ParseMigration(string(content))
	if err != nil {
		return err
	}

	return nil
}
