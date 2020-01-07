package app

import (
	"fmt"
	"github.com/Knetic/govaluate"
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
	if params.ColorSpace == "" || params.C1 == "" || params.C2 == "" || params.C3 == "" {
		return colorful.Hsv(c, c, c)
	}

	var color colorful.Color

	parameters := make(map[string]interface{}, 8)
	parameters["c"] = c

	var C1 *govaluate.EvaluableExpression
	var C2 *govaluate.EvaluableExpression
	var C3 *govaluate.EvaluableExpression
	functions := map[string]govaluate.ExpressionFunction{
		"tanh": func(args ...interface{}) (interface{}, error) {
			return math.Tanh(args[0].(float64)), nil
		},
	}

	exp1, err1 := govaluate.NewEvaluableExpressionWithFunctions(params.C1, functions)
	if err1 != nil {
		panic(fmt.Sprintf("First param error: %s", err1))
	} else {
		C1 = exp1
	}
	exp2, err2 := govaluate.NewEvaluableExpressionWithFunctions(params.C2, functions)
	if err2 != nil {
		panic(fmt.Sprintf("First param error: %s", err2))
	} else {
		C2 = exp2
	}
	exp3, err3 := govaluate.NewEvaluableExpressionWithFunctions(params.C2, functions)
	if err3 != nil {
		panic(fmt.Sprintf("First param error: %s", err3))
	} else {
		C3 = exp3
	}

	c1, err1 := C1.Evaluate(parameters)
	c2, err2 := C2.Evaluate(parameters)
	c3, err3 := C3.Evaluate(parameters)

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
