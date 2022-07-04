package testing

import "github.com/ory/dockertest/v3"

type DatabaseConfig interface {
	DockerOptions() *dockertest.RunOptions
	DatabaseUrl(resource *dockertest.Resource) string
}
