package main

import (
	"fmt"
	"os"
)

type GenerateParams struct {
	imageSize      float32
	c              complex64
	generationMode string
	outputFile     string
	centerX        float32
	centerY        float32
	axisSpan       float32
}

// smooth coloring: https://iquilezles.org/www/articles/mset_smooth/mset_smooth.htm

func main() {
	var params GenerateParams
	// init application environment
	initEnv()
	if len(os.Args) > 5 {
		// parse cl arguments
		params = GenerateParams{
			imageSize:      toFloat(os.Args[1]),
			c:              complex(toFloat(os.Args[2]), toFloat(os.Args[3])),
			generationMode: os.Args[4],
			outputFile:     fmt.Sprintf("out/%s", os.Args[5]),
			centerX:        0,
			centerY:        0,
			axisSpan:       2,
		}
		printPrams(os.Args)
	} else {
		// set static params for development
		params = GenerateParams{
			imageSize:      1000.0,
			c:              complex(-0.6, 0),
			generationMode: "-i",
			outputFile:     "out/out-test.png",
			centerX:        -1.415,
			centerY:        0,
			axisSpan:       0.02,
		}
	}
	generateImage(params)
	fmt.Println("DONE !")
}

func initEnv() {
	// init file outputs target directory
	exists, _ := exists("../out")
	if exists {
		return
	}
	err := os.Mkdir("../out", os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func generateImage(params GenerateParams) {
	calcParams := SetParams{
		maxY:           params.centerY + params.axisSpan,
		minY:           params.centerY - params.axisSpan,
		minX:           params.centerX - params.axisSpan,
		maxX:           params.centerX + params.axisSpan,
		step:           params.axisSpan / (params.imageSize * 2),
		maxIterations:  10,
		thresholdValue: 40,
		c:              params.c,
	}
	if params.generationMode == "-i" {
		set := CalcByIterations(calcParams)
		DrawByIteration(set, int(params.imageSize), int(params.imageSize), params.outputFile)
	} else if params.generationMode == "-t" {
		set := CalcByThreshold(calcParams)
		DrawByThreshold(set, int(params.imageSize), int(params.imageSize), params.outputFile)
	} else {
		panic(fmt.Sprintf("Invalid mode %s", params.generationMode))
	}
}
