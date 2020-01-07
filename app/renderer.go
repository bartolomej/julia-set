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
	"strings"
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
	if renderParams.RenderMode == "ITERATION" {
		set := CalcByIterations(setParams)
		img := renderImgByIteration(set, renderParams)
		_ = saveImage(renderParams.Filename, img, renderParams.Encoding, renderParams.Id)
	} else if renderParams.RenderMode == "THRESHOLD" {
		set := CalcByThreshold(setParams)
		img := renderImgByThreshold(set, renderParams)
		_ = saveImage(renderParams.Filename, img, renderParams.Encoding, renderParams.Id)
	} else {
		panic(fmt.Sprintf("Invalid RenderMode %s", renderParams.RenderMode))
	}
}

//c := array[y*stepY][x*stepX] / 80
//color := colorful.Hsv(c * 60 + 100, 0, c) // black and white
//color := colorful.Hsv(c * 60 + 240, 1, math.Tanh(c * 2)) // purple
//color := colorful.Hsv(c * 60, 1, math.Tanh(c * 2)) // yellow

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
			c := array[y*stepY][x*stepX]
			im.Set(x, y, evalColor(c, params.Color))
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
			im.Set(x, y, evalColor(c, params.Color))
		}
	}
	return im
}

func evalColor(c float64, params ColorParams) colorful.Color {
	// return default color when params empty
	if params.ColorSpace == "" || params.C1 == nil || params.C2 == nil || params.C3 == nil {
		return colorful.Hsv(c, c, c)
	}

	var color colorful.Color

	parameters := make(map[string]interface{}, 8)
	parameters["c"] = c

	c1, err1 := params.C1.Evaluate(parameters)
	c2, err2 := params.C1.Evaluate(parameters)
	c3, err3 := params.C1.Evaluate(parameters)

	if err1 != nil {
		panic(fmt.Sprintf("First color param eval error: %s", err1))
	}
	if err2 != nil {
		panic(fmt.Sprintf("Second color param eval error: %s", err2))
	}
	if err3 != nil {
		panic(fmt.Sprintf("Third color param eval error: %s", err3))
	}

	if strings.EqualFold(params.ColorSpace, "HSV") {
		color = colorful.Hsv(c1.(float64), c2.(float64), c3.(float64))
	} else if strings.EqualFold(params.ColorSpace, "RGB") {
		color = colorful.LinearRgb(c1.(float64), c2.(float64), c3.(float64))
	} else {
		panic(fmt.Sprintf("Color space %s not supported", params.ColorSpace))
	}

	return color
}

func saveImage(filename string, im image.Image, encoding string, id string) error {
	// make folder for current configuration
	folder := fmt.Sprintf("out/%s", id)
	MakeDir(folder)
	file, err := os.Create(fmt.Sprintf("%s/%s.%s", folder, filename, encoding))
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
