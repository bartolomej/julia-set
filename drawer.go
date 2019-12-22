package main

import (
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"image/png"
	"math"
	"math/cmplx"
	"os"
)

func DrawByThreshold(array [][]float64, width int, height int, file string) {
	im := image.NewNRGBA64(image.Rect(0, 0, width, height))
	if len(array) < height || len(array[0]) < width {
		panic("Array smaller than drawing size")
	}
	stepY := len(array) / height
	stepX := len(array[0]) / width
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := array[y*stepY][x*stepX] * 20
			color := colorful.Hsv(c, c, c)
			im.Set(x, y, color)
		}
	}
	_ = savePNG(file, im)
}

func DrawByIteration(array [][]complex128, width int, height int, file string) {
	im := image.NewNRGBA64(image.Rect(0, 0, width, height))
	if len(array) < height || len(array[0]) < width {
		panic("Array smaller than drawing size")
	}
	stepY := len(array) / height
	stepX := len(array[0]) / width
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := math.Pow(math.E, -cmplx.Abs(array[y*stepY][x*stepX]))
			color := colorful.Hsv(c, c, c)
			im.Set(x, y, color)
		}
	}
	_ = savePNG(file, im)
}

func savePNG(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, im)
}
