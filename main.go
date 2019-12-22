package main

import (
	"encoding/json"
	"fmt"
	"github.com/bartolomej/complex-set-art/app"
	"os"
)

type GenerateParams struct {
	imageSize      float32
	outputFile     string
	c              complex64
	generationMode string
	centerX        float32
	centerY        float32
	axisSpan       float32
}

// smooth coloring: https://iquilezles.org/www/articles/mset_smooth/mset_smooth.htm

func main() {
	var params GenerateParams
	// init out folder if doesn't exist
	app.MakeDir("../out")
	if len(os.Args) == 2 {
		params = parseFileParams(os.Args[1])
	} else if len(os.Args) > 4 {
		// parse cl arguments
		params = parseCliParams()
	} else {
		// set static params for development
		params = GenerateParams{
			imageSize:      1000.0,
			c:              complex(-0.3, 0),
			generationMode: "-t",
			outputFile:     "out/out-test.png",
			centerX:        0,
			centerY:        0,
			axisSpan:       2,
		}
	}
	// print currently used params
	printParams(params)
	generateImage(params)
	fmt.Println("DONE !")
}

func parseFileParams(id string) GenerateParams {
	str, err := app.ReadFile("renders.json")
	if err != nil {
		panic("Error opening renders.json config file")
	}
	var res []interface{}
	_ = json.Unmarshal([]byte(str), &res)
	for i := 0; i < len(res); i++ {
		var obj map[string]interface{}
		data, _ := json.Marshal(res[i])
		_ = json.Unmarshal(data, &obj)
		configId := obj["id"].(string)
		realC := obj["realC"].(float64)
		imagC := obj["imagC"].(float64)
		mode := obj["mode"].(string)
		if configId == id {
			return GenerateParams{
				imageSize:      float32(obj["imageSize"].(float64)),
				outputFile:     fmt.Sprintf("out/r%fi%f_%s.png", realC, imagC, mode),
				c:              complex64(complex(obj["realC"].(float64), obj["imagC"].(float64))),
				generationMode: obj["mode"].(string),
				centerX:        float32(obj["centerX"].(float64)),
				centerY:        float32(obj["centerY"].(float64)),
				axisSpan:       float32(obj["axisSpan"].(float64)),
			}
		}
	}
	panic(fmt.Sprintf("Config with id: %s not found", id))
}

func parseCliParams() GenerateParams {
	return GenerateParams{
		imageSize:      app.ParamToFloat(os.Args[1]),
		c:              complex(app.ParamToFloat(os.Args[2]), app.ParamToFloat(os.Args[3])),
		generationMode: os.Args[4],
		outputFile:     fmt.Sprintf("out/r%si%s_%s.png", os.Args[2], os.Args[3], os.Args[4]),
		centerX:        0,
		centerY:        0,
		axisSpan:       2,
	}
}

func printParams(params GenerateParams) {
	fmt.Printf("Image size: %f \n", params.imageSize)
	fmt.Printf("Hyperparam C: %f + %fi \n", real(params.c), imag(params.c))
	fmt.Printf("Generation mode: %s \n", params.generationMode)
	fmt.Printf("Output file: %s \n", params.outputFile)
}

func generateImage(params GenerateParams) {
	calcParams := app.SetParams{
		MaxY:           params.centerY + params.axisSpan,
		MinY:           params.centerY - params.axisSpan,
		MinX:           params.centerX - params.axisSpan,
		MaxX:           params.centerX + params.axisSpan,
		Step:           params.axisSpan / (params.imageSize * 2),
		MaxIterations:  30,
		ThresholdValue: 40,
		C:              params.c,
	}
	if params.generationMode == "-i" {
		set := app.CalcByIterations(calcParams)
		app.DrawByIteration(set, int(params.imageSize), int(params.imageSize), params.outputFile)
	} else if params.generationMode == "-t" {
		set := app.CalcByThreshold(calcParams)
		app.DrawByThreshold(set, int(params.imageSize), int(params.imageSize), params.outputFile)
	} else {
		panic(fmt.Sprintf("Invalid mode %s", params.generationMode))
	}
}
