package yamlsorter_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(Fail)
	RunSpecs(m, "YAML Sorter Suite")
}
