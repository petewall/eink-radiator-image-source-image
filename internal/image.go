package internal

import (
	"image/color"
	"io"

	"github.com/fogleman/gg"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . ImageContext
type ImageContext interface {
	DrawRectangle(x, y, w, h float64)
	Fill()
	SetColor(c color.Color)
	EncodePNG(w io.Writer) error
	SavePNG(path string) error
}

//counterfeiter:generate . ImageContextMaker
type ImageContextMaker func(width, height int) ImageContext

var NewImageContext ImageContextMaker = func(width, height int) ImageContext {
	return gg.NewContext(width, height)
}
