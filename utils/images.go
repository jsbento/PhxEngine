package utils

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/davecgh/go-spew/spew"
)

type ImageFormat string

const (
	ImageJPEG ImageFormat = "jpeg"
	ImagePNG  ImageFormat = "png"
)

func (f ImageFormat) Valid() error {
	switch f {
	case ImageJPEG, ImagePNG:
		return nil
	default:
		return errors.New("invalid image format")
	}
}

func LoadImage(path string) (image.Image, ImageInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, ImageInfo{}, err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, ImageInfo{}, err
	}
	if err := ImageFormat(format).Valid(); err != nil {
		return nil, ImageInfo{}, err
	}
	spew.Dump(img)
	spew.Dump(format)
	return img, GetImageInfo(img, ImageFormat(format)), err
}

type ImageInfo struct {
	Width    int
	Height   int
	Channels int
	Format   ImageFormat
}

func GetImageInfo(img image.Image, format ImageFormat) ImageInfo {
	bounds := img.Bounds()

	return ImageInfo{
		Width:    bounds.Dx(),
		Height:   bounds.Dy(),
		Channels: getNumChannels(img),
		Format:   format,
	}
}

func EncodeJPEG(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func EncodePNG(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getNumChannels(img image.Image) int {
	switch img.(type) {
	case *image.Gray, *image.Gray16:
		return 1
	case *image.YCbCr:
		return 3
	case *image.RGBA, *image.NRGBA, *image.RGBA64, *image.NRGBA64:
		return 4
	case *image.CMYK:
		return 4
	case *image.Paletted:
		return 1
	default:
		return 4
	}
}
