package mysql

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMySqlDriver(t *testing.T) {
	_t = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "MySql Driver Test Suite")
}

var _t *testing.T
