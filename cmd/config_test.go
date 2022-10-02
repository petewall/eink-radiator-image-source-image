package cmd_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/petewall/eink-radiator-image-source-image/v2/cmd"
)

var _ = Describe("Config", func() {
	var output *Buffer

	BeforeEach(func() {
		output = NewBuffer()
		cmd.ConfigCmd.SetOut(output)
	})

	It("prints a blank config", func() {
		cmd.ConfigCmd.Run(cmd.ConfigCmd, []string{})
		Expect(output).Should(Say("{\"color\":\"\"}"))
	})
})
