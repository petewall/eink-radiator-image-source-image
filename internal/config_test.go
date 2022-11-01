package internal_test

import (
	"encoding/json"
	"errors"
	"image"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	"github.com/petewall/eink-radiator-image-source-image/internal"
	"github.com/petewall/eink-radiator-image-source-image/internal/internalfakes"
)

var _ = Describe("Config", func() {
	Describe("GenerateImage", func() {
		var (
			httpGetter   *internalfakes.FakeHttpGetter
			imageDecoder *internalfakes.FakeImageDecoder
		)

		BeforeEach(func() {
			img := image.NewRGBA(image.Rect(0, 0, 1024, 768))
			imageDecoder = &internalfakes.FakeImageDecoder{}
			imageDecoder.Returns(img, nil)

			res := &http.Response{}
			httpGetter = &internalfakes.FakeHttpGetter{}
			httpGetter.Returns(res, nil)

			internal.DecodeImage = imageDecoder.Spy
			internal.HttpGet = httpGetter.Spy
		})

		Context("resized image", func() {
			It("fetches an image and returns a scaled image", func() {
				config := &internal.Config{
					Source: "https://www.example.com/link.jpg",
					Scale:  "resize",
					Backgound: &internal.BackgroundType{
						Color: "red",
					},
				}

				img, err := config.GenerateImage(200, 300)
				Expect(err).ToNot(HaveOccurred())

				By("fetching the image", func() {
					Expect(httpGetter.CallCount()).To(Equal(1))
					Expect(httpGetter.ArgsForCall(0)).To(Equal("https://www.example.com/link.jpg"))
				})

				By("decoding the image", func() {
					Expect(imageDecoder.CallCount()).To(Equal(1))
				})

				By("returning a resized version of the image", func() {
					Expect(img.Bounds().Max).To(Equal(image.Point{200, 300}))
				})
			})
		})

		Context("contained image", func() {
			It("fetches an image and returns a contained image", func() {
				config := &internal.Config{
					Source: "https://www.example.com/link.jpg",
					Scale:  "contain",
					Backgound: &internal.BackgroundType{
						Color: "red",
					},
				}

				img, err := config.GenerateImage(200, 300)
				Expect(err).ToNot(HaveOccurred())

				By("fetching the image", func() {
					Expect(httpGetter.CallCount()).To(Equal(1))
					Expect(httpGetter.ArgsForCall(0)).To(Equal("https://www.example.com/link.jpg"))
				})

				By("decoding the image", func() {
					Expect(imageDecoder.CallCount()).To(Equal(1))
				})

				By("returning a contained version of the image", func() {
					Expect(img.Bounds().Max).To(Equal(image.Point{200, 300}))
				})
			})
		})

		Context("covered image", func() {
			It("fetches an image and returns a covered image", func() {
				config := &internal.Config{
					Source: "https://www.example.com/link.jpg",
					Scale:  "cover",
					Backgound: &internal.BackgroundType{
						Color: "red",
					},
				}

				img, err := config.GenerateImage(200, 300)
				Expect(err).ToNot(HaveOccurred())

				By("fetching the image", func() {
					Expect(httpGetter.CallCount()).To(Equal(1))
					Expect(httpGetter.ArgsForCall(0)).To(Equal("https://www.example.com/link.jpg"))
				})

				By("decoding the image", func() {
					Expect(imageDecoder.CallCount()).To(Equal(1))
				})

				By("returning a covered version of the image", func() {
					Expect(img.Bounds().Max).To(Equal(image.Point{200, 300}))
				})
			})
		})

		Context("unknown scale type", func() {
			It("returns an error", func() {
				config := &internal.Config{
					Source: "https://www.example.com/link.jpg",
					Scale:  "smoothjazz",
					Backgound: &internal.BackgroundType{
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
				config := &internal.Config{
					Source: "https://www.example.com/link.jpg",
					Scale:  "cover",
					Backgound: &internal.BackgroundType{
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
				config := &internal.Config{
					Source: "https://www.example.com/link.jpg",
					Scale:  "cover",
					Backgound: &internal.BackgroundType{
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
		config := internal.Config{
			Source: "https://www.example.com/link.jpg",
			Scale:  "contain",
			Backgound: &internal.BackgroundType{
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
		config, err := internal.ParseConfig(configFile.Name())
		Expect(err).ToNot(HaveOccurred())
		Expect(config.Source).To(Equal("https://www.example.com/link.jpg"))
		Expect(config.Scale).To(Equal("contain"))
		Expect(config.Backgound.Color).To(Equal("red"))
	})

	Context("config file is json formatted", func() {
		BeforeEach(func() {
			config := internal.Config{
				Source: "https://www.example.com/impa.jpg",
				Scale:  "cover",
				Backgound: &internal.BackgroundType{
					Color: "blue",
				},
			}
			var err error
			configFileContents, err = json.Marshal(config)
			Expect(err).ToNot(HaveOccurred())
		})

		It("parses just fine", func() {
			config, err := internal.ParseConfig(configFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(config.Source).To(Equal("https://www.example.com/impa.jpg"))
			Expect(config.Scale).To(Equal("cover"))
			Expect(config.Backgound.Color).To(Equal("blue"))
		})
	})

	When("reading the config file fails", func() {
		It("returns an error", func() {
			_, err := internal.ParseConfig("this file does not exist")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to read image config file: open this file does not exist: no such file or directory"))
		})
	})

	When("parsing the config file fails", func() {
		BeforeEach(func() {
			configFileContents = []byte("this is invalid yaml!")
		})

		It("returns an error", func() {
			_, err := internal.ParseConfig(configFile.Name())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to parse image config file: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `this is...` into internal.Config"))
		})
	})

	When("the config file has invalid data", func() {
		BeforeEach(func() {
			config := internal.Config{
				Source: "https://www.example.com/impa.jpg",
				Scale:  "zelda",
				Backgound: &internal.BackgroundType{
					Color: "link",
				},
			}
			var err error
			configFileContents, err = json.Marshal(config)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error", func() {
			_, err := internal.ParseConfig(configFile.Name())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("config file is not valid: scale value is invalid: \"zelda\", must be one of resize, contain, cover"))
		})
	})
})
