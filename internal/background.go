package internal

import (
	"image"

	blank "github.com/petewall/eink-radiator-image-source-blank/pkg"
)

//counterfeiter:generate . BackgroundMaker
type BackgroundMaker func(width, height int, color string) image.Image

var MakeBackground BackgroundMaker = func(width, height int, color string) image.Image {
	backgroundConfig := blank.Config{Color: color}
	return backgroundConfig.GenerateImage(width, height)
}
