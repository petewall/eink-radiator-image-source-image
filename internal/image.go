package internal

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"os"

	"golang.org/x/image/draw"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . ImageDecoder
type ImageDecoder func(r io.Reader) (image.Image, error)

var DecodeImage ImageDecoder = func(r io.Reader) (image.Image, error) {
	im, _, err := image.Decode(r)
	return im, err
}

//counterfeiter:generate . ImageEncoder
type ImageEncoder func(w io.Writer, i image.Image) error

var EncodeImage ImageEncoder = png.Encode

//counterfeiter:generate . ImageWriter
type ImageWriter func(file string, i image.Image) error

var WriteImage ImageWriter = func(file string, i image.Image) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	err = EncodeImage(f, i)
	if err != nil {
		return err
	}
	return f.Close()
}

//counterfeiter:generate . ImageMaker
type ImageMaker func(width, height int) *image.RGBA

var NewImage ImageMaker = func(width, height int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, width, height))
}

//counterfeiter:generate . ImageScaler
type ImageScaler func(dst draw.Image, dr image.Rectangle, src image.Image, sr image.Rectangle, op draw.Op, opts *draw.Options)

var Scale ImageScaler = draw.CatmullRom.Scale

//counterfeiter:generate . Drawer
type Drawer func(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, op draw.Op)

var Draw Drawer = draw.Draw
