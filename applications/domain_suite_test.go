package applications

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestApplication(t *testing.T) {
	_t = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Application Test Suite")
}

var _t *testing.T
