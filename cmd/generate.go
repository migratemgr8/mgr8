package cmd

import (
	"fmt"
	"os"

	"github.com/kenji-yamane/mgr8/domain"
)

type generate struct{}

func (g *generate) execute(args []string, databaseURL string, driver domain.Driver) error {
	fileNameOld := args[0]
	content, err := os.ReadFile(fileNameOld)
	if err != nil {
		return fmt.Errorf("could not read from old file: %s", err)
	}

	oldSchema, err := driver.ParseMigration(string(content))
	if err != nil {
		return fmt.Errorf("could not parse migration from old file: %s", err)
	}

	fileNameNew := args[1]
	content, err = os.ReadFile(fileNameNew)
	if err != nil {
		return fmt.Errorf("could not read from new file: %s", err)
	}

	newSchema, err := driver.ParseMigration(string(content))
	if err != nil {
		return fmt.Errorf("could not parse migration from new file: %s", err)
	}

	diffQueue := newSchema.Diff(oldSchema)
	deparser := driver.Deparser()
	for _, diff := range diffQueue {
		up := diff.Up(deparser)
		println(up)
	}

	return nil
}
