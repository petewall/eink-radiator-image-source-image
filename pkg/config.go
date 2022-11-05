package pkg

import (
	"fmt"
	"image"
	"math"
	"os"

	"golang.org/x/image/draw"

	"gopkg.in/yaml.v2"

	blank "github.com/petewall/eink-radiator-image-source-blank/pkg"
	"github.com/petewall/eink-radiator-image-source-image/internal"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . ImageGenerator
type ImageGenerator interface {
	GenerateImage(width, height int) (image.Image, error)
}

const (
	ScaleResize  = "resize"
	ScaleContain = "contain"
	ScaleCover   = "cover"
)

type BackgroundType struct {
	Color string `json:"color" yaml:"color"`
}

type Config struct {
	Source     string          `json:"source" yaml:"source"`
	Scale      string          `json:"scale" yaml:"scale"`
	Background *BackgroundType `json:"background,omitempty" yaml:"background,omitempty"`
}

func (c *Config) GenerateImage(width, height int) (image.Image, error) {
	res, err := internal.HttpGet(c.Source)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image (%s): %w", c.Source, err)
	}

	im, err := internal.DecodeImage(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image (%s): %w", c.Source, err)
	}

	switch c.Scale {
	case ScaleContain:
		return c.generateContainedImage(width, height, im)
	case ScaleCover:
		return c.generateCoveredImage(width, height, im)
	case ScaleResize:
		return c.generateResizedImage(width, height, im)
	default:
		return nil, fmt.Errorf("unknown image scale type: %s", c.Scale)
	}
}

func (c *Config) generateContainedImage(width, height int, im image.Image) (image.Image, error) {
	background := internal.MakeBackground(width, height, c.Background.Color)

	xScale := float64(width) / float64(im.Bounds().Size().X)
	yScale := float64(height) / float64(im.Bounds().Size().Y)
	scaleFactor := math.Min(xScale, yScale)

	scaledWidth := int(scaleFactor * float64(im.Bounds().Size().X))
	scaledHeight := int(scaleFactor * float64(im.Bounds().Size().Y))
	scaled := internal.NewImage(scaledWidth, scaledHeight)
	internal.Scale(scaled, scaled.Rect, im, im.Bounds(), draw.Over, nil)

	var sp image.Point
	if xScale > yScale {
		sp = image.Point{X: (scaledWidth - width) / 2, Y: 0}
	} else {
		sp = image.Point{X: 0, Y: (scaledHeight - height) / 2}
	}

	dst := internal.NewImage(width, height)
	internal.Draw(dst, dst.Rect, background, image.Point{}, draw.Src)
	internal.Draw(dst, dst.Rect, scaled, sp, draw.Src)
	return dst, nil
}

func (c *Config) generateCoveredImage(width, height int, im image.Image) (image.Image, error) {
	xScale := float64(width) / float64(im.Bounds().Size().X)
	yScale := float64(height) / float64(im.Bounds().Size().Y)
	scaleFactor := math.Max(xScale, yScale)

	scaledWidth := int(scaleFactor * float64(im.Bounds().Size().X))
	scaledHeight := int(scaleFactor * float64(im.Bounds().Size().Y))
	scaled := internal.NewImage(scaledWidth, scaledHeight)
	internal.Scale(scaled, scaled.Rect, im, im.Bounds(), draw.Over, nil)

	var sp image.Point
	if xScale > yScale {
		sp = image.Point{X: 0, Y: (scaledHeight - height) / 2}
	} else {
		sp = image.Point{X: (scaledWidth - width) / 2, Y: 0}
	}

	dst := internal.NewImage(width, height)
	internal.Draw(dst, dst.Rect, scaled, sp, draw.Over)
	return dst, nil
}

func (c *Config) generateResizedImage(width, height int, im image.Image) (image.Image, error) {
	dst := internal.NewImage(width, height)
	internal.Scale(dst, dst.Rect, im, im.Bounds(), draw.Over, nil)
	return dst, nil
}

func (c *Config) Validate() error {
	if c.Source == "" {
		return fmt.Errorf("missing image source")
	}

	if c.Scale != ScaleResize &&
		c.Scale != ScaleContain &&
		c.Scale != ScaleCover {
		return fmt.Errorf("scale value is invalid: \"%s\", must be one of resize, contain, cover", c.Scale)
	}

	backgroundConfig := blank.Config{Color: c.Background.Color}
	if err := backgroundConfig.Validate(); err != nil {
		return fmt.Errorf("invalid background: %w", err)
	}
	return nil
}

func ParseConfig(path string) (*Config, error) {
	configData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read image config file: %w", err)
	}

	var config *Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image config file: %w", err)
	}

	if config.Background == nil {
		config.Background = &BackgroundType{
			Color: "white",
		}
	}
	if config.Background.Color == "" {
		config.Background.Color = "white"
	}

	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("config file is not valid: %w", err)
	}

	return config, nil
}
