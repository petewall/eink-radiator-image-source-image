package test_test

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Generate", func() {
	var outputFile string

	AfterEach(func() {
		if outputFile != "" {
			Expect(os.Remove(outputFile)).To(Succeed())
			outputFile = ""
		}
	})

	Context("using resize", func() {
		It("generates an image that is resized to fit the desired resolution", func() {
			outputFile = "resized.png"
			Run("generate --config resize-config.yaml --output resized.png --height 300 --width 400")
			Eventually(CommandSession).WithTimeout(time.Second * 5).Should(Exit(0))

			By("saving the image to a file", func() {
				actualData, err := os.ReadFile(outputFile)
				Expect(err).ToNot(HaveOccurred())
				expectedData, err := os.ReadFile("expected_resized.png")
				Expect(err).ToNot(HaveOccurred())
				Expect(actualData).To(Equal(expectedData))
			})
		})
	})

	Context("using contain", func() {
		It("generates an image that is shrunk to fit the desired resolution", func() {
			outputFile = "contained.png"
			Run("generate --config contain-config.json --output contained.png --height 300 --width 300")
			Eventually(CommandSession).WithTimeout(time.Second * 5).Should(Exit(0))

			By("saving the image to a file", func() {
				actualData, err := os.ReadFile(outputFile)
				Expect(err).ToNot(HaveOccurred())
				expectedData, err := os.ReadFile("expected_contained.png")
				Expect(err).ToNot(HaveOccurred())
				Expect(actualData).To(Equal(expectedData))
			})
		})
	})

	Context("using cover", func() {
		It("generates an image that is cropped to fit the desired resolution", func() {
			outputFile = "covered.png"
			Run("generate --config cover-config.yaml --output covered.png --height 300 --width 400")
			Eventually(CommandSession).WithTimeout(time.Second * 5).Should(Exit(0))

			By("saving the image to a file", func() {
				actualData, err := os.ReadFile(outputFile)
				Expect(err).ToNot(HaveOccurred())
				expectedData, err := os.ReadFile("expected_covered.png")
				Expect(err).ToNot(HaveOccurred())
				Expect(actualData).To(Equal(expectedData))
			})
		})
	})

	When("using --to-stdout", func() {
		It("writes the image to stdout", func() {
			Run("generate --config resize-config.yaml --to-stdout --height 300 --width 400")
			Eventually(CommandSession).WithTimeout(time.Second * 5).Should(Exit(0))

			expectedData, err := os.ReadFile("expected_resized.png")
			Expect(err).ToNot(HaveOccurred())
			Expect(CommandSession.Out.Contents()).To(Equal(expectedData))
		})
	})
})
