package main

import (
	"github.com/lucasb-eyer/go-colorful"
	"image/png"
	"strconv"

	"fmt"
	"image"
	"math"
	"math/cmplx"
	"os"
)

type SetParams struct {
	maxY       float32
	minY       float32
	minX       float32
	maxX       float32
	step       float32
	iterations int
	c          complex128
}

func main() {
	if len(os.Args) != 5 {
		panic("Invalid number of args")
	}
	size := toFloat(os.Args[1])
	realC := toFloat(os.Args[2])
	imagC := toFloat(os.Args[3])
	file := os.Args[4]

	fmt.Println(size, realC, imagC, file)
	params := SetParams{
		maxY:       2,
		minY:       -2,
		minX:       -2,
		maxX:       2,
		step:       2 / size,
		iterations: 10,
		c:          complex128(complex(realC, imagC)),
	}
	set := calcSet(params)
	draw(set, int(size)*2, int(size)*2, file)
}

func toFloat(input string) float32 {
	n, err := strconv.ParseFloat(input, 32)
	if err != nil {
		panic(fmt.Sprintf("Parameter %s is not a number", input))
	}
	return float32(n)
}

func calcSet(params SetParams) [][]complex128 {
	var minY float32 = -2
	var maxY float32 = 2
	var minX float32 = -2
	var maxX float32 = 2
	var y float32 = 0
	var x float32 = 0
	var results [][]complex128
	for y = maxY; y >= minY; y -= params.step {
		var resX []complex128
		for x = minX; x <= maxX; x += params.step {
			z := complex128(complex(x, y))
			for i := 0; i < params.iterations; i++ {
				z = cmplx.Pow(z, 2) + params.c
			}
			resX = append(resX, z)
		}
		results = append(results, resX)
	}
	return results
}

func draw(array [][]complex128, width int, height int, file string) {
	im := image.NewNRGBA64(image.Rect(0, 0, width, height))
	if len(array) < height || len(array[0]) < width {
		panic("Array smaller than drawing size")
	}

	stepY := len(array) / height
	stepX := len(array[0]) / width
	for y := 0; y < height; y += stepY {
		for x := 0; x < width; x += stepX {
			c := math.Pow(math.E, -cmplx.Abs(array[y][x]))
			color := colorful.Hsv(c, c, c)
			im.Set(x, y, color)
		}
	}

	_ = savePNG(file, im)
}

func savePNG(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, im)
}
