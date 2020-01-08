package app

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Render(params RenderParams) {
	params.Folder = "out"
	if (params.Video != VideoParams{}) {
		renderVideo(params)
	} else if (params.Image != AbstractParams{}) {
		renderImage(params)
	}
}

func renderVideo(params RenderParams) {
	totalFrames := int(float64(params.Video.Fps) * params.Video.Duration)

	centerX := params.Video.Start.OriginX
	diffX := params.Video.End.OriginX - params.Video.Start.OriginX
	stepX := diffX / float32(totalFrames)

	centerY := params.Video.Start.OriginY
	diffY := params.Video.End.OriginY - params.Video.Start.OriginY
	stepY := diffY / float32(totalFrames)

	axisSpan := params.Video.Start.AxisSpan
	diffSpan := params.Video.End.AxisSpan - params.Video.Start.AxisSpan
	stepSpan := diffSpan / float32(totalFrames)

	C := params.Video.Start.C
	diffC := params.Video.End.C - params.Video.Start.C
	stepC := diffC / complex(float32(totalFrames), float32(totalFrames))

	exp := params.Video.Start.Exponent
	diffExp := params.Video.End.Exponent - params.Video.Start.Exponent
	stepExp := diffExp / complex(float32(totalFrames), float32(totalFrames))

	digits := strconv.Itoa(len(strconv.Itoa(totalFrames)))

	for i := 0; i < totalFrames; i++ {
		params.Filename = fmt.Sprintf("frame%0"+digits+"d", i)
		params.Image = AbstractParams{
			OriginX:  centerX,
			OriginY:  centerY,
			AxisSpan: axisSpan,
			Exponent: exp,
			C:        C,
		}
		renderImage(params)
		centerX += stepX
		centerY += stepY
		axisSpan += stepSpan
		exp += stepExp
		C += stepC
	}

	inputDir := "cache/" + params.Id + "/frame%0" + digits + "d.png"
	outputDir := fmt.Sprintf("out/%s.mp4", params.Id)

	// DOCS: https://trac.ffmpeg.org/wiki/Slideshow
	// ffmpeg -framerate 30 -i frame%02d.png video.mp4
	cmd := exec.Command(
		"ffmpeg",
		"-framerate",
		strconv.Itoa(params.Video.Fps),
		"-i",
		inputDir,
		outputDir,
	)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(stdout))
}

func renderImage(params RenderParams) {
	setParams := SetParams{
		OriginX:      params.Image.OriginX,
		OriginY:      params.Image.OriginY,
		Resolution:   params.Resolution,
		AxisSpan:     params.Image.AxisSpan,
		MaxDistance:  params.MaxThreshold,
		MaxIteration: params.MaxIterations,
		ReturnMode:   params.ReturnMode,
		Exponent:     params.Image.Exponent,
		C:            params.Image.C,
	}
	set := CalculateSet(setParams)
	size := int(params.Resolution)
	img := image.NewNRGBA64(image.Rect(0, 0, size, size))
	stepY := len(set) / size
	stepX := len(set[0]) / size
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			c := set[y*stepY][x*stepX]
			// TODO: SMOOTH ITERATION: c := math.Pow(math.E, -set[y*stepY][x*stepX])
			img.Set(x, y, evalColor(c, params.Color))
		}
		progress := (100 * y) / size
		fmt.Print("\rRENDERING IMAGE: " + strconv.Itoa(progress) + "%")
	}
	fmt.Println()
	saveImage(params.Folder, params.Filename, img, params.Encoding, params.Id)
}

func evalColor(c float64, params ColorParams) colorful.Color {
	// return default color when params empty
	if params.ColorSpace == "" || params.C1 == "" || params.C2 == "" || params.C3 == "" {
		return colorful.Hsv(c, c, c)
	}

	var color colorful.Color

	parameters := make(map[string]interface{}, 8)
	parameters["c"] = c

	functions := map[string]govaluate.ExpressionFunction{
		"tanh": func(args ...interface{}) (interface{}, error) {
			return math.Tanh(args[0].(float64)), nil
		},
		"tan": func(args ...interface{}) (interface{}, error) {
			return math.Tan(args[0].(float64)), nil
		},
	}

	var C1 *govaluate.EvaluableExpression
	var C2 *govaluate.EvaluableExpression
	var C3 *govaluate.EvaluableExpression

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

func saveImage(folder string, filename string, im image.Image, encoding string, id string) {
	// make folder for current configuration
	path := fmt.Sprintf("%s/%s", folder, id)
	MakeDir(path)
	file, err := os.Create(fmt.Sprintf("%s/%s.%s", path, filename, encoding))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if encoding == "png" {
		err := png.Encode(file, im)
		if err != nil {
			panic(err)
		}
	} else if encoding == "jpeg" {
		err := jpeg.Encode(file, im, nil)
		if err != nil {
			panic(err)
		}
	} else {
		panic(fmt.Sprintf("Invalid encoding: %s", encoding))
	}
}
