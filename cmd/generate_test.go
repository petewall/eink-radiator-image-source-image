package cmd_test

import (
	"errors"
	"image"

	"github.com/spf13/viper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/petewall/eink-radiator-image-source-image/cmd"
	"github.com/petewall/eink-radiator-image-source-image/internal"
	"github.com/petewall/eink-radiator-image-source-image/internal/internalfakes"
	"github.com/petewall/eink-radiator-image-source-image/pkg/pkgfakes"
)

var _ = Describe("Generate", func() {
	var (
		img            image.Image
		imageGenerator *pkgfakes.FakeImageGenerator
		imageEncoder   *internalfakes.FakeImageEncoder
		imageWriter    *internalfakes.FakeImageWriter
	)

	BeforeEach(func() {
		img = image.NewRGBA(image.Rect(0, 0, 10, 10))
		imageGenerator = &pkgfakes.FakeImageGenerator{}
		imageGenerator.GenerateImageReturns(img, nil)
		imageEncoder = &internalfakes.FakeImageEncoder{}
		imageWriter = &internalfakes.FakeImageWriter{}

		cmd.ImageGenerator = imageGenerator
		internal.EncodeImage = imageEncoder.Spy
		internal.WriteImage = imageWriter.Spy

		viper.Set("to-stdout", false)
		viper.Set("output", cmd.DefaultOutputFilename)
		viper.Set("height", 1000)
		viper.Set("width", 2000)
	})

	It("generates a blank image", func() {
		err := cmd.GenerateCmd.RunE(cmd.GenerateCmd, []string{})
		Expect(err).ToNot(HaveOccurred())

		By("calling the image generator", func() {
			Expect(imageGenerator.GenerateImageCallCount()).To(Equal(1))
		})

		By("defaulting to writing to image.png", func() {
			Expect(imageEncoder.CallCount()).To(Equal(0))
			Expect(imageWriter.CallCount()).To(Equal(1))
			filename, writtenImage := imageWriter.ArgsForCall(0)
			Expect(filename).To(Equal("image.png"))
			Expect(writtenImage).To(Equal(img))
		})

		By("using the right resolution", func() {
			Expect(imageGenerator.GenerateImageCallCount()).To(Equal(1))
			width, height := imageGenerator.GenerateImageArgsForCall(0)
			Expect(width).To(Equal(2000))
			Expect(height).To(Equal(1000))
		})
	})

	When("using --to-stdout", func() {
		var output *Buffer

		BeforeEach(func() {
			output = NewBuffer()
			cmd.GenerateCmd.SetOut(output)
			viper.Set("to-stdout", true)
		})

		It("outputs the image to stdout", func() {
			err := cmd.GenerateCmd.RunE(cmd.GenerateCmd, []string{})
			Expect(err).ToNot(HaveOccurred())

			Expect(imageWriter.CallCount()).To(Equal(0))
			Expect(imageEncoder.CallCount()).To(Equal(1))
			out, encodedImage := imageEncoder.ArgsForCall(0)
			Expect(out).To(Equal(output))
			Expect(encodedImage).To(Equal(img))
		})

		When("encoding fails", func() {
			BeforeEach(func() {
				imageEncoder.Returns(errors.New("encode image failed"))
			})

			It("returns an error", func() {
				err := cmd.GenerateCmd.RunE(cmd.GenerateCmd, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("encode image failed"))
			})
		})
	})

	When("saving the image fails", func() {
		BeforeEach(func() {
			imageWriter.Returns(errors.New("save image failed"))
		})

		It("returns an error", func() {
			err := cmd.GenerateCmd.RunE(cmd.GenerateCmd, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("save image failed"))
		})
	})
})
