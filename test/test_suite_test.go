package test_test

import (
	"os/exec"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}

var (
	imageSourcePath string
	CommandSession  *gexec.Session
)

var _ = BeforeSuite(func() {
	var err error
	imageSourcePath, err = gexec.Build(
		"github.com/petewall/eink-radiator-image-source-image/v2",
		"-ldflags",
		"-X github.com/petewall/eink-radiator-image-source-image/v2/cmd.Version=1.2.3",
	)
	Expect(err).NotTo(HaveOccurred())
})

func Run(args string) {
	var err error
	command := exec.Command(imageSourcePath, strings.Split(args, " ")...)
	CommandSession, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())
}

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
	gexec.KillAndWait()
})
