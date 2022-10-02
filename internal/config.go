package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	ScaleResize  = "resize"
	ScaleContain = "contain"
	ScaleCover   = "cover"
)

type BackgroundType struct {
	Color string `json:"color" yaml:"color"`
}

type Config struct {
	Source    string          `json:"src" yaml:"src"`
	Scale     string          `json:"scale" yaml:"scale"`
	Backgound *BackgroundType `json:"background,omitempty" yaml:"background,omitempty"`
}

func (c *Config) GenerateImage(width, height int) ImageContext {
	imageContext := NewImageContext(width, height)
	// imageContext.SetColor(c.GetColor())
	imageContext.DrawRectangle(0, 0, float64(width), float64(height))
	imageContext.Fill()

	return imageContext
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
