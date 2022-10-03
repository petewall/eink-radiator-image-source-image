package internal

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"os"
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
