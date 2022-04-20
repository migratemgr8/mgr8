package postgres

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPostgresDriver(t *testing.T) {
	_t = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Postgres Driver Test Suite")
}

var _t *testing.T
