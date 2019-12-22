package app

import (
	"math/cmplx"
)

type SetParams struct {
	MaxY           float32
	MinY           float32
	MinX           float32
	MaxX           float32
	Step           float32
	MaxIterations  int
	ThresholdValue int
	C              complex64
}

func CalcByIterations(params SetParams) [][]complex128 {
	var y float32 = 0
	var x float32 = 0
	var results [][]complex128
	for y = params.MaxY; y >= params.MinY; y -= params.Step {
		var resX []complex128
		for x = params.MinX; x <= params.MaxX; x += params.Step {
			z := complex128(complex(x, y))
			for i := 0; i < params.MaxIterations; i++ {
				z = cmplx.Pow(z, 2) + complex128(params.C)
			}
			resX = append(resX, z)
		}
		results = append(results, resX)
	}
	return results
}

func CalcByThreshold(params SetParams) [][]float64 {
	var y float32 = 0
	var x float32 = 0
	var results [][]float64
	for y = params.MaxY; y >= params.MinY; y -= params.Step {
		var resX []float64
		for x = params.MinX; x <= params.MaxX; x += params.Step {
			z := complex128(complex(x, y))
			i := 0
			//var r float64
			for cmplx.Abs(z) < float64(params.ThresholdValue) && i < params.MaxIterations {
				z = cmplx.Pow(z, 2) + complex128(params.C)
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
