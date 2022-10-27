package internal

import (
	"fmt"
	"image"
	"net/http"
	"os"

	"golang.org/x/image/draw"

	"gopkg.in/yaml.v2"
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
	Source    string          `json:"source" yaml:"source"`
	Scale     string          `json:"scale" yaml:"scale"`
	Backgound *BackgroundType `json:"background,omitempty" yaml:"background,omitempty"`
}

//counterfeiter:generate . HttpGetter
type HttpGetter func(path string) (*http.Response, error)

var HttpGet = http.Get

func (c *Config) GenerateImage(width, height int) (image.Image, error) {
	res, err := HttpGet(c.Source)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image (%s): %w", c.Source, err)
	}

	im, err := DecodeImage(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image (%s): %w", c.Source, err)
	}

	switch c.Scale {
	case ScaleContain:
		return c.generateContainedImage(width, height, im)
	case ScaleCover:
		return c.generateCoveredImage(width, height, im)
	case ScaleResize:
		return c.generateScaledImage(width, height, im)
	default:
		return nil, fmt.Errorf("unknown image scale type: %s", c.Scale)
	}
}

func (c *Config) generateContainedImage(width, height int, im image.Image) (image.Image, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (c *Config) generateCoveredImage(width, height int, im image.Image) (image.Image, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (c *Config) generateScaledImage(width, height int, im image.Image) (image.Image, error) {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.ApproxBiLinear.Scale(dst, dst.Rect, im, im.Bounds(), draw.Over, nil)
	return dst, nil
}

func (c *Config) Validate() error {
	if c.Scale != ScaleResize &&
		c.Scale != ScaleContain &&
		c.Scale != ScaleCover {
		return fmt.Errorf("scale value is invalid: \"%s\", must be one of resize, contain, cover", c.Scale)
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

	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("config file is not valid: %w", err)
	}

	return config, nil
}
