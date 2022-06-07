package domain_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDomain(t *testing.T) {
	_t = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Domain Test Suite")
}

var _t *testing.T
