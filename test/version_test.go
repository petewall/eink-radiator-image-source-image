package test_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Version", func() {
	It("returns the version number", func() {
		Run("version")
		Eventually(CommandSession).Should(Exit(0))
		Eventually(CommandSession.Out).Should(Say("eInk Radiator image image source"))
		Eventually(CommandSession.Out).Should(Say("Version: 1.2.3"))
	})
})
