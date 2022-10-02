package cmd_test

import (
	"errors"

	"github.com/spf13/viper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/petewall/eink-radiator-image-source-image/v2/cmd"
	"github.com/petewall/eink-radiator-image-source-image/v2/internal"
	"github.com/petewall/eink-radiator-image-source-image/v2/internal/internalfakes"
)

var _ = Describe("Generate", func() {
	var (
		newImageContext *internalfakes.FakeImageContextMaker
		imageContext    *internalfakes.FakeImageContext
	)

	BeforeEach(func() {
		imageContext = &internalfakes.FakeImageContext{}
		newImageContext = &internalfakes.FakeImageContextMaker{}
		newImageContext.Returns(imageContext)
		internal.NewImageContext = newImageContext.Spy

		viper.Set("to-stdout", false)
		viper.Set("output", cmd.DefaultOutputFilename)
		viper.Set("height", cmd.DefaultImageHeight)
		viper.Set("width", cmd.DefaultImageWidth)
	})

	It("generates a blank image", func() {
		cmd.Config = &internal.Config{
			Source: "https://www.example.com/ganon.jpg",
			Scale:  "cover",
			Backgound: &internal.BackgroundType{
				Color: "red",
			},
		}
		err := cmd.GenerateCmd.RunE(cmd.GenerateCmd, []string{})
		Expect(err).ToNot(HaveOccurred())

		By("defaulting to writing to image.png", func() {
			Expect(imageContext.SavePNGCallCount()).To(Equal(1))
			Expect(imageContext.SavePNGArgsForCall(0)).To(Equal("image.png"))
		})

		By("defaulting to 640x480", func() {
			Expect(newImageContext.CallCount()).To(Equal(1))
			width, height := newImageContext.ArgsForCall(0)
			Expect(width).To(Equal(640))
			Expect(height).To(Equal(480))
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
			cmd.Config = &internal.Config{
				Source: "https://www.example.com/moblin.jpg",
				Scale:  "cover",
				Backgound: &internal.BackgroundType{
					Color: "black",
				},
			}
			err := cmd.GenerateCmd.RunE(cmd.GenerateCmd, []string{})
			Expect(err).ToNot(HaveOccurred())

			Expect(imageContext.EncodePNGCallCount()).To(Equal(1))
			Expect(imageContext.EncodePNGArgsForCall(0)).To(Equal(output))
		})

		When("encoding fails", func() {
			BeforeEach(func() {
				imageContext.EncodePNGReturns(errors.New("encode png failed"))
			})

			It("returns an error", func() {
				err := cmd.GenerateCmd.RunE(cmd.GenerateCmd, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("encode png failed"))
			})
		})
	})

	When("using --height and --width to change the resolution", func() {
		It("generates an image of the specified resolution", func() {
			viper.Set("height", 1000)
			viper.Set("width", 2000)

			cmd.Config = &internal.Config{
				Source: "https://www.example.com/lizafos.jpg",
				Scale:  "cover",
				Backgound: &internal.BackgroundType{
					Color: "green",
				},
			}
			err := cmd.GenerateCmd.RunE(cmd.GenerateCmd, []string{})
			Expect(err).ToNot(HaveOccurred())

			Expect(newImageContext.CallCount()).To(Equal(1))
			width, height := newImageContext.ArgsForCall(0)
			Expect(width).To(Equal(2000))
			Expect(height).To(Equal(1000))
		})
	})

	When("saving the image fails", func() {
		BeforeEach(func() {
			imageContext.SavePNGReturns(errors.New("save png failed"))
		})

		It("returns an error", func() {
			err := cmd.GenerateCmd.RunE(cmd.GenerateCmd, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("save png failed"))
		})
	})
})
