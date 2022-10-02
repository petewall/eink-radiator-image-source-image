package test_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"gopkg.in/yaml.v2"

	"github.com/petewall/eink-radiator-image-source-image/v2/internal"
)

var _ = Describe("Config", func() {
	It("returns the blank config", func() {
		Run("config")
		Eventually(CommandSession).Should(Exit(0))
		output := CommandSession.Out.Contents()
		var blankConfig internal.Config
		Expect(yaml.Unmarshal(output, &blankConfig)).To(Succeed())
		Expect(blankConfig.Source).To(BeEmpty())
		Expect(blankConfig.Scale).To(Equal("resize"))
		Expect(blankConfig.Backgound.Color).To(BeEmpty())
	})
})
