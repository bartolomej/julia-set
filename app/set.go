package app

import (
	"fmt"
	"math/cmplx"
)

type SetParams struct {
	CenterX       float32
	CenterY       float32
	Resolution    float32
	AxisSpan      float32
	MaxThreshold  int
	MaxIterations int
	C             complex64
}

type calcParams struct {
	maxY           float32
	minY           float32
	minX           float32
	maxX           float32
	step           float32
	maxIterations  int
	thresholdValue int
	c              complex64
}

func getParams(params SetParams) calcParams {
	return calcParams{
		maxY:           params.CenterY + params.AxisSpan,
		minY:           params.CenterY - params.AxisSpan,
		minX:           params.CenterX - params.AxisSpan,
		maxX:           params.CenterX + params.AxisSpan,
		step:           params.AxisSpan / (params.Resolution * 2),
		maxIterations:  params.MaxIterations,
		thresholdValue: params.MaxThreshold,
		c:              params.C,
	}
}

func CalcByIterations(set SetParams) [][]complex128 {
	params := getParams(set)
	var y float32 = 0
	var x float32 = 0
	var results [][]complex128
	fmt.Println("COMPUTING: ")
	total := (params.maxY - params.minY) / params.step
	for y = params.maxY; y >= params.minY; y -= params.step {
		var resX []complex128
		for x = params.minX; x <= params.maxX; x += params.step {
			z := complex128(complex(x, y))
			for i := 0; i < params.maxIterations; i++ {
				z = cmplx.Pow(z, 2) + complex128(params.c)
			}
			resX = append(resX, z)
		}
		results = append(results, resX)
		// compute progress
		current := (y - params.minY) / params.step
		progress := int(float64(100 - (100*current)/total))
		fmt.Printf("\r%d %", progress)
	}
	return results
}

func CalcByThreshold(set SetParams) [][]float64 {
	params := getParams(set)
	var y float32 = 0
	var x float32 = 0
	var results [][]float64
	fmt.Println("COMPUTING: ")
	total := (params.maxY - params.minY) / params.step
	for y = params.maxY; y >= params.minY; y -= params.step {
		var resX []float64
		for x = params.minX; x <= params.maxX; x += params.step {
			z := complex128(complex(x, y))
			i := 0
			for cmplx.Abs(z) < float64(params.thresholdValue) && i < params.maxIterations {
				z = cmplx.Pow(z, 2) + complex128(params.c)
				i++
			}
			resX = append(resX, float64(i))
		}
		results = append(results, resX)
		// compute progress
		current := (y - params.minY) / params.step
		progress := int(float64(100 - (100*current)/total))
		fmt.Printf("\r%d %", progress)
	}
	return results
}
