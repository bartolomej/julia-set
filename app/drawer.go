package app

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"math/cmplx"
	"os"
)

func DrawByThreshold(array [][]float64, params RenderParams) {
	size := int(params.Resolution)
	im := image.NewNRGBA64(image.Rect(0, 0, size, size))
	if len(array) < size || len(array[0]) < size {
		panic("Array smaller than drawing size")
	}
	stepY := len(array) / size
	stepX := len(array[0]) / size
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			c := array[y*stepY][x*stepX] * 25
			color := colorful.LinearRgb(c, c/2, c/3)
			im.Set(x, y, color)
		}
	}
	_ = saveImage(params.Filename, im, params.Encoding)
}

func DrawByIteration(array [][]complex128, params RenderParams) {
	size := int(params.Resolution)
	im := image.NewNRGBA64(image.Rect(0, 0, size, size))
	if len(array) < size || len(array[0]) < size {
		panic("Array smaller than drawing size")
	}
	stepY := len(array) / size
	stepX := len(array[0]) / size
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			c := math.Pow(math.E, -cmplx.Abs(array[y*stepY][x*stepX]))
			color := colorful.Hsv(c, c, c)
			im.Set(x, y, color)
		}
	}
	_ = saveImage(params.Filename, im, params.Encoding)
}

func saveImage(filename string, im image.Image, encoding string) error {
	file, err := os.Create(fmt.Sprintf("out/%s.%s", filename, encoding))
	if err != nil {
		return err
	}
	defer file.Close()
	if encoding == "png" {
		return png.Encode(file, im)
	} else if encoding == "jpeg" {
		return jpeg.Encode(file, im, nil)
	} else {
		panic(fmt.Sprintf("Invalid encoding: %s", encoding))
	}
}
