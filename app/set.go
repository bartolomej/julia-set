package app

import (
	"fmt"
	"math"
	"math/cmplx"
	"strconv"
)

type SetParams struct {
	OriginX      float32
	OriginY      float32
	Resolution   float32
	AxisSpan     float32
	MaxDistance  int
	MaxIteration int
	ReturnMode   ReturnMode
	Exponent     complex64
	C            complex64
}

type ReturnMode string

const (
	ITERATION ReturnMode = "ITERATION"
	DISTANCE  ReturnMode = "DISTANCE"
)

func CalculateSet(set SetParams) [][]float64 {
	var y float32 = 0
	var x float32 = 0
	var results [][]float64

	maxY := set.OriginY + set.AxisSpan
	minY := set.OriginY - set.AxisSpan
	maxX := set.OriginX + set.AxisSpan
	minX := set.OriginX - set.AxisSpan
	step := set.AxisSpan / (set.Resolution * 2)
	totalSteps := (maxY - minY) / step

	for y = maxY; y >= minY; y -= step {
		var resX []float64
		for x = minX; x <= maxX; x += step {
			z := complex128(complex(x, y))
			i := 0
			if set.MaxDistance != 0 && set.MaxIteration != 0 {
				for cmplx.Abs(z) < float64(set.MaxDistance) && i < set.MaxIteration {
					z = cmplx.Pow(z, complex128(set.Exponent)) + complex128(set.C)
					i++
				}
			} else if set.MaxIteration != 0 {
				for i < set.MaxIteration {
					z = cmplx.Pow(z, complex128(set.Exponent)) + complex128(set.C)
					i++
				}
			} else if set.MaxDistance != 0 {
				for cmplx.Abs(z) < float64(set.MaxDistance) {
					z = cmplx.Pow(z, complex128(set.Exponent)) + complex128(set.C)
					i++
				}
			}
			if set.ReturnMode == ITERATION {
				resX = append(resX, smoothIter(i, cmplx.Abs(complex128(set.Exponent)), z))
			} else if set.ReturnMode == DISTANCE {
				resX = append(resX, cmplx.Abs(z))
			} else {
				panic("Invalid return mode: " + set.ReturnMode)
			}
		}
		results = append(results, resX)

		// compute progress
		current := (y - minY) / step
		progress := int(float64(100 - (100*current)/totalSteps))
		fmt.Print("\rCOMPUTING: " + strconv.Itoa(progress) + "%")
	}
	fmt.Println()
	return results
}

// SMOOTH ITERATION COUNT TAKEN FROM:
// http://linas.org/art-gallery/escape/smooth.html
func smoothIter(n int, d float64, z complex128) float64 {
	return float64(n) + 1 - math.Log(math.Log(cmplx.Abs(z)))/math.Log(d)
}
