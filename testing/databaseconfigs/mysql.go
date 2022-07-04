package databaseconfigs

import (
	"fmt"
	"github.com/ory/dockertest/v3"
)

type mySqlConfig struct {
	repository string
	tag        string
	userName   string
	password   string
	dbName     string
	port       string
}

func NewMySqlConfig() *mySqlConfig {
	return &mySqlConfig{
		repository: "mysql",
		tag:        "5.7",
		userName:   "root",
		password:   "root",
		dbName:     "core",
		port:       "3306/tcp",
	}
}

func (c *mySqlConfig) DockerOptions() *dockertest.RunOptions {
	return &dockertest.RunOptions{
		Repository: c.repository,
		Tag:        c.tag,
		Env: []string{
			fmt.Sprintf("MYSQL_ROOT_PASSWORD=%v", c.password),
			fmt.Sprintf("MYSQL_DATABASE=%v", c.dbName),
		},
	}
}

func (c *mySqlConfig) DatabaseUrl(resource *dockertest.Resource) string {
	hostPort := resource.GetPort(c.port)
	return fmt.Sprintf("%s:%s@(localhost:%s)/%s?parseTime=true", c.userName, c.password, hostPort, c.dbName)
}
