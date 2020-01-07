package main

import (
	"fmt"
	"github.com/bartolomej/complex-set-art/app"
	"os"
	"strconv"
)

func main() {
	var params app.RenderParams

	// init app folders
	app.MakeDir("out")
	app.MakeDir("cache")

	if len(os.Args) == 2 && os.Args[1] == "default-image" {
		params = getDefaultImageParams()
	} else if len(os.Args) == 2 && os.Args[1] == "default-video" {
		panic("Video rendering not implemented yet")
	} else if len(os.Args) == 2 {
		params = app.ParseFileParams(os.Args[1], "")
	} else if len(os.Args) == 3 {
		params = app.ParseFileParams(os.Args[1], os.Args[2])
	} else if len(os.Args) > 4 {
		params = parseCliParams()
	} else {
		params = app.ParseFileParams("video-test", "")
	}
	app.Render(params)
	fmt.Println("\n --> DONE !")
}

func parseCliParams() app.RenderParams {
	res := paramToFloat(os.Args[1])
	c := complex(paramToFloat(os.Args[2]), paramToFloat(os.Args[3]))
	filename := fmt.Sprintf("r%fi%f_%s", real(c), imag(c), os.Args[4])
	// use static config for cli args
	image := app.AbstractParams{
		C:        c,
		CenterX:  0,
		CenterY:  0,
		AxisSpan: 2,
	}
	return app.RenderParams{
		Resolution: res,
		RenderMode: os.Args[4],
		Encoding:   "png",
		Filename:   filename,
		Image:      image,
	}
}

func getDefaultImageParams() app.RenderParams {
	return app.RenderParams{
		Id:            "default",
		Resolution:    100,
		RenderMode:    "THRESHOLD",
		Encoding:      "png",
		Filename:      "test-out",
		MaxThreshold:  30,
		MaxIterations: 20,
		Image: app.AbstractParams{
			C:        complex(0, 0),
			CenterX:  0,
			CenterY:  0,
			AxisSpan: 2,
		},
	}
}

func paramToFloat(input string) float32 {
	n, err := strconv.ParseFloat(input, 32)
	if err != nil {
		panic(fmt.Sprintf("Parameter %s is not a number", input))
	}
	return float32(n)
}
