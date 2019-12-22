package app

import (
	"math/cmplx"
)

type SetParams struct {
	CenterX    float32
	CenterY    float32
	Resolution float32
	AxisSpan   float32
	C          complex64
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
		maxIterations:  30,
		thresholdValue: 40,
		c:              params.C,
	}
}

func CalcByIterations(set SetParams) [][]complex128 {
	params := getParams(set)
	var y float32 = 0
	var x float32 = 0
	var results [][]complex128
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
	}
	return results
}

func CalcByThreshold(set SetParams) [][]float64 {
	params := getParams(set)
	var y float32 = 0
	var x float32 = 0
	var results [][]float64
	for y = params.maxY; y >= params.minY; y -= params.step {
		var resX []float64
		for x = params.minX; x <= params.maxX; x += params.step {
			z := complex128(complex(x, y))
			i := 0
			//var r float64
			for cmplx.Abs(z) < float64(params.thresholdValue) && i < params.maxIterations {
				z = cmplx.Pow(z, 2) + complex128(params.c)
				// TODO: implement smooth iteration coloring
				// https://en.wikibooks.org/wiki/Fractals/Iterations_in_the_complex_plane/Julia_set
				//r = float64(i) - math.Log2(math.Log2(cmplx.Abs(z)))
				// https://iquilezles.org/www/articles/mset_smooth/mset_smooth.htm
				i++
			}
			resX = append(resX, float64(i))
		}
		results = append(results, resX)
	}
	return results
}
