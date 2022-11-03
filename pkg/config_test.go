package pkg_test

import (
	"encoding/json"
	"errors"
	"image"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/image/draw"
	"gopkg.in/yaml.v2"

	"github.com/petewall/eink-radiator-image-source-image/internal"
	"github.com/petewall/eink-radiator-image-source-image/internal/internalfakes"
	"github.com/petewall/eink-radiator-image-source-image/pkg"
)

var _ = Describe("Config", func() {
	Describe("GenerateImage", func() {
		var (
			httpGetter   *internalfakes.FakeHttpGetter
			imageDecoder *internalfakes.FakeImageDecoder

			backgroundImage *image.RGBA
			fetchedImage    *image.RGBA
			returnedImage   *image.RGBA

			makeBackground *internalfakes.FakeBackgroundMaker
			drawer         *internalfakes.FakeDrawer
			newImage       *internalfakes.FakeImageMaker
			scale          *internalfakes.FakeImageScaler
		)

		BeforeEach(func() {
			fetchedImage = image.NewRGBA(image.Rect(0, 0, 1024, 768))
			imageDecoder = &internalfakes.FakeImageDecoder{}
			imageDecoder.Returns(fetchedImage, nil)
			internal.DecodeImage = imageDecoder.Spy

			res := &http.Response{}
			httpGetter = &internalfakes.FakeHttpGetter{}
			httpGetter.Returns(res, nil)
			internal.HttpGet = httpGetter.Spy

			backgroundImage = image.NewRGBA(image.Rect(0, 0, 300, 200))
			makeBackground = &internalfakes.FakeBackgroundMaker{}
			makeBackground.Returns(backgroundImage)
			internal.MakeBackground = makeBackground.Spy

			drawer = &internalfakes.FakeDrawer{}
			internal.Draw = drawer.Spy

			returnedImage = image.NewRGBA(image.Rect(0, 0, 300, 200))
			newImage = &internalfakes.FakeImageMaker{}
			newImage.ReturnsOnCall(0, returnedImage)
			internal.NewImage = newImage.Spy

			scale = &internalfakes.FakeImageScaler{}
			internal.Scale = scale.Spy
		})

		Context("resized image", func() {
			It("fetches an image and returns a scaled image", func() {
				config := &pkg.Config{
					Source: "https://www.example.com/link.jpg",
					Scale:  "resize",
					Background: &pkg.BackgroundType{
						Color: "red",
					},
				}

				img, err := config.GenerateImage(300, 200)
				Expect(err).ToNot(HaveOccurred())

				By("fetching the image", func() {
					Expect(httpGetter.CallCount()).To(Equal(1))
					Expect(httpGetter.ArgsForCall(0)).To(Equal("https://www.example.com/link.jpg"))
				})

				By("decoding the image", func() {
					Expect(imageDecoder.CallCount()).To(Equal(1))
				})

				By("returning a resized version of the image", func() {
					Expect(img).To(Equal(returnedImage))

					Expect(newImage.CallCount()).To(Equal(1))
					width, height := newImage.ArgsForCall(0)
					Expect(width).To(Equal(300))
					Expect(height).To(Equal(200))

					Expect(scale.CallCount()).To(Equal(1))
					dst, rect, im, imRect, op, options := scale.ArgsForCall(0)
					Expect(dst).To(Equal(returnedImage))
					Expect(rect).To(Equal(image.Rect(0, 0, 300, 200)))
					Expect(im).To(Equal(fetchedImage))
					Expect(imRect).To(Equal(image.Rect(0, 0, 1024, 768)))
					Expect(op).To(Equal(draw.Over))
					Expect(options).To(BeNil())
				})
			})
		})

		Context("contained image", func() {
			var scaledImage *image.RGBA
			BeforeEach(func() {
				scaledImage = image.NewRGBA(image.Rect(0, 0, 266, 200))
				newImage.ReturnsOnCall(0, scaledImage)
				newImage.ReturnsOnCall(1, returnedImage)
			})

			It("fetches an image and returns a covered image", func() {
				config := &pkg.Config{
					Source: "https://www.example.com/link.jpg",
					Scale:  "contain",
					Background: &pkg.BackgroundType{
						Color: "red",
					},
				}

				img, err := config.GenerateImage(300, 200)
				Expect(err).ToNot(HaveOccurred())

				By("fetching the image", func() {
					Expect(httpGetter.CallCount()).To(Equal(1))
					Expect(httpGetter.ArgsForCall(0)).To(Equal("https://www.example.com/link.jpg"))
				})

				By("decoding the image", func() {
					Expect(imageDecoder.CallCount()).To(Equal(1))
				})

				By("building a background", func() {
					Expect(makeBackground.CallCount()).To(Equal(1))
					width, height, color := makeBackground.ArgsForCall(0)
					Expect(width).To(Equal(300))
					Expect(height).To(Equal(200))
					Expect(color).To(Equal("red"))
				})

				By("scaling the image to the right size", func() {
					Expect(newImage.CallCount()).To(Equal(2))
					width, height := newImage.ArgsForCall(0)
					Expect(width).To(Equal(266))
					Expect(height).To(Equal(200))

					Expect(scale.CallCount()).To(Equal(1))
					dst, dstRect, im, imRect, op, options := scale.ArgsForCall(0)
					Expect(dst).To(Equal(scaledImage))
					Expect(dstRect).To(Equal(image.Rect(0, 0, 266, 200)))
					Expect(im).To(Equal(fetchedImage))
					Expect(imRect).To(Equal(image.Rect(0, 0, 1024, 768)))
					Expect(op).To(Equal(draw.Over))
					Expect(options).To(BeNil())
				})

				By("drawing the background onto the returned image", func() {
					Expect(drawer.CallCount()).To(Equal(2))
					dst, dstRect, scaled, scaledPoint, op := drawer.ArgsForCall(0)
					Expect(dst).To(Equal(returnedImage))
					Expect(dstRect).To(Equal(image.Rect(0, 0, 300, 200)))
					Expect(scaled).To(Equal(backgroundImage))
					Expect(scaledPoint).To(Equal(image.Point{0, 0}))
					Expect(op).To(Equal(draw.Src))
				})

				By("drawing the scaled image onto the returned image", func() {
					dst, dstRect, scaled, scaledPoint, op := drawer.ArgsForCall(1)
					Expect(dst).To(Equal(returnedImage))
					Expect(dstRect).To(Equal(image.Rect(0, 0, 300, 200)))
					Expect(scaled).To(Equal(scaledImage))
					Expect(scaledPoint).To(Equal(image.Point{-17, 0}))
					Expect(op).To(Equal(draw.Src))
				})

				By("returning a contained version of the image", func() {
					Expect(img).To(Equal(returnedImage))

					width, height := newImage.ArgsForCall(1)
					Expect(width).To(Equal(300))
					Expect(height).To(Equal(200))
				})
			})
		})

		Context("covered image", func() {
			var scaledImage *image.RGBA
			BeforeEach(func() {
				scaledImage = image.NewRGBA(image.Rect(0, 0, 300, 225))
				newImage.ReturnsOnCall(0, scaledImage)
				newImage.ReturnsOnCall(1, returnedImage)
			})

			It("fetches an image and returns a covered image", func() {
				config := &pkg.Config{
					Source: "https://www.example.com/link.jpg",
					Scale:  "cover",
					Background: &pkg.BackgroundType{
						Color: "red",
					},
				}

				img, err := config.GenerateImage(300, 200)
				Expect(err).ToNot(HaveOccurred())

				By("fetching the image", func() {
					Expect(httpGetter.CallCount()).To(Equal(1))
					Expect(httpGetter.ArgsForCall(0)).To(Equal("https://www.example.com/link.jpg"))
				})

				By("decoding the image", func() {
					Expect(imageDecoder.CallCount()).To(Equal(1))
				})

				By("scaling the image to the right size", func() {
					Expect(newImage.CallCount()).To(Equal(2))
					width, height := newImage.ArgsForCall(0)
					Expect(width).To(Equal(300))
					Expect(height).To(Equal(225))

					Expect(scale.CallCount()).To(Equal(1))
					dst, dstRect, im, imRect, op, options := scale.ArgsForCall(0)
					Expect(dst).To(Equal(scaledImage))
					Expect(dstRect).To(Equal(image.Rect(0, 0, 300, 225)))
					Expect(im).To(Equal(fetchedImage))
					Expect(imRect).To(Equal(image.Rect(0, 0, 1024, 768)))
					Expect(op).To(Equal(draw.Over))
					Expect(options).To(BeNil())
				})

				By("drawing the scaled image onto the returned image", func() {
					Expect(drawer.CallCount()).To(Equal(1))
					dst, dstRect, scaled, scaledPoint, op := drawer.ArgsForCall(0)
					Expect(dst).To(Equal(returnedImage))
					Expect(dstRect).To(Equal(image.Rect(0, 0, 300, 200)))
					Expect(scaled).To(Equal(scaledImage))
					Expect(scaledPoint).To(Equal(image.Point{0, 12}))
					Expect(op).To(Equal(draw.Over))
				})

				By("returning a covered version of the image", func() {
					Expect(img).To(Equal(returnedImage))

					width, height := newImage.ArgsForCall(1)
					Expect(width).To(Equal(300))
					Expect(height).To(Equal(200))
				})
			})
		})

		Context("unknown scale type", func() {
			It("returns an error", func() {
				config := &pkg.Config{
					Source: "https://www.example.com/link.jpg",
					Scale:  "smoothjazz",
					Background: &pkg.BackgroundType{
						Color: "red",
					},
				}

				_, err := config.GenerateImage(200, 300)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("unknown image scale type: smoothjazz"))
			})
		})

		When("getting the image fails", func() {
			BeforeEach(func() {
				httpGetter.Returns(nil, errors.New("http get failed"))
			})

			It("returns an error", func() {
				config := &pkg.Config{
					Source: "https://www.example.com/link.jpg",
					Scale:  "cover",
					Background: &pkg.BackgroundType{
						Color: "red",
					},
				}

				_, err := config.GenerateImage(200, 300)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to fetch image (https://www.example.com/link.jpg): http get failed"))
			})
		})

		When("decoding the image fails", func() {
			BeforeEach(func() {
				imageDecoder.Returns(nil, errors.New("image decoding failed"))
			})

			It("returns an error", func() {
				config := &pkg.Config{
					Source: "https://www.example.com/link.jpg",
					Scale:  "cover",
					Background: &pkg.BackgroundType{
						Color: "red",
					},
				}

				_, err := config.GenerateImage(200, 300)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to decode image (https://www.example.com/link.jpg): image decoding failed"))
			})
		})
	})
})

var _ = Describe("ParseConfig", func() {
	var (
		configFile         *os.File
		configFileContents []byte
	)

	JustBeforeEach(func() {
		var err error
		configFile, err = os.CreateTemp("", "blank-config.yaml")
		Expect(err).ToNot(HaveOccurred())
		_, err = configFile.Write(configFileContents)
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		config := pkg.Config{
			Source: "https://www.example.com/link.jpg",
			Scale:  "contain",
			Background: &pkg.BackgroundType{
				Color: "red",
			},
		}
		var err error
		configFileContents, err = yaml.Marshal(config)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.Remove(configFile.Name())).To(Succeed())
	})

	It("parses the image config file", func() {
		config, err := pkg.ParseConfig(configFile.Name())
		Expect(err).ToNot(HaveOccurred())
		Expect(config.Source).To(Equal("https://www.example.com/link.jpg"))
		Expect(config.Scale).To(Equal("contain"))
		Expect(config.Background.Color).To(Equal("red"))
	})

	Context("config file is json formatted", func() {
		BeforeEach(func() {
			config := pkg.Config{
				Source: "https://www.example.com/impa.jpg",
				Scale:  "cover",
				Background: &pkg.BackgroundType{
					Color: "blue",
				},
			}
			var err error
			configFileContents, err = json.Marshal(config)
			Expect(err).ToNot(HaveOccurred())
		})

		It("parses just fine", func() {
			config, err := pkg.ParseConfig(configFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(config.Source).To(Equal("https://www.example.com/impa.jpg"))
			Expect(config.Scale).To(Equal("cover"))
			Expect(config.Background.Color).To(Equal("blue"))
		})
	})

	When("reading the config file fails", func() {
		It("returns an error", func() {
			_, err := pkg.ParseConfig("this file does not exist")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to read image config file: open this file does not exist: no such file or directory"))
		})
	})

	When("parsing the config file fails", func() {
		BeforeEach(func() {
			configFileContents = []byte("this is invalid yaml!")
		})

		It("returns an error", func() {
			_, err := pkg.ParseConfig(configFile.Name())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to parse image config file: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `this is...` into pkg.Config"))
		})
	})

	When("the config file has invalid data", func() {
		BeforeEach(func() {
			config := pkg.Config{
				Source: "https://www.example.com/impa.jpg",
				Scale:  "zelda",
				Background: &pkg.BackgroundType{
					Color: "link",
				},
			}
			var err error
			configFileContents, err = json.Marshal(config)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error", func() {
			_, err := pkg.ParseConfig(configFile.Name())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("config file is not valid: scale value is invalid: \"zelda\", must be one of resize, contain, cover"))
		})
	})
})
