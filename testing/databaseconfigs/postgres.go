package databaseconfigs

import (
	"fmt"
	"github.com/ory/dockertest/v3"
)

type postgresConfig struct {
	repository string
	tag        string
	userName   string
	password   string
	dbName     string
	port       string
}

func NewPostgresConfig() *postgresConfig {
	return &postgresConfig{
		repository: "postgres",
		tag:        "11",
		userName:   "root",
		password:   "root",
		dbName:     "core",
		port:       "5432/tcp",
	}
}

func (c *postgresConfig) DockerOptions() *dockertest.RunOptions {
	return &dockertest.RunOptions{
		Repository: c.repository,
		Tag:        c.tag,
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%v", c.password),
			fmt.Sprintf("POSTGRES_USER=%v", c.userName),
			fmt.Sprintf("POSTGRES_DB=%v", c.dbName),
			"listen_addresses = '*'",
		},
	}
}

func (c *postgresConfig) DatabaseUrl(resource *dockertest.Resource) string {
	hostPort := resource.GetHostPort(c.port)
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.userName, c.password, hostPort, c.dbName)
}
