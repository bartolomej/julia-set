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

func RenderImage(renderParams RenderParams) {
	setParams := SetParams{
		CenterX:       renderParams.Image.CenterX,
		CenterY:       renderParams.Image.CenterY,
		Resolution:    renderParams.Resolution,
		AxisSpan:      renderParams.Image.AxisSpan,
		C:             renderParams.Image.C,
		MaxThreshold:  renderParams.MaxThreshold,
		MaxIterations: renderParams.MaxIterations,
	}
	if renderParams.RenderMode == "-i" {
		set := CalcByIterations(setParams)
		img := renderImgByIteration(set, renderParams)
		_ = saveImage(renderParams.Filename, img, renderParams.Encoding)
	} else if renderParams.RenderMode == "-t" {
		set := CalcByThreshold(setParams)
		img := renderImgByThreshold(set, renderParams)
		_ = saveImage(renderParams.Filename, img, renderParams.Encoding)
	} else {
		panic(fmt.Sprintf("Invalid RenderMode %s", renderParams.RenderMode))
	}
}

func renderImgByThreshold(array [][]float64, params RenderParams) *image.NRGBA64 {
	size := int(params.Resolution)
	im := image.NewNRGBA64(image.Rect(0, 0, size, size))
	if len(array) < size || len(array[0]) < size {
		panic("Array smaller than drawing size")
	}
	stepY := len(array) / size
	stepX := len(array[0]) / size
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			//c := array[y*stepY][x*stepX] / 80
			c := array[y*stepY][x*stepX]
			//color := colorful.Hsv(c * 60 + 100, 0, c) // black and white
			//color := colorful.Hsv(c * 60 + 240, 1, math.Tanh(c * 2)) // purple
			//color := colorful.Hsv(c * 60, 1, math.Tanh(c * 2)) // yellow
			color := colorful.LinearRgb(c, c/2, c)
			im.Set(x, y, color)
		}
	}
	return im
}

func renderImgByIteration(array [][]complex128, params RenderParams) *image.NRGBA64 {
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
	return im
}

func saveImage(filename string, im image.Image, encoding string) error {
	file, err := os.Create(fmt.Sprintf("%s.%s", filename, encoding))
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
