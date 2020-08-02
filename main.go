package main

import (
	"fmt"
	"github.com/bartolomej/complex-set-art/app"
	"os"
	"strconv"
)

const (
	ExampleConfig   = "example.config.json"
	Config          = "config.json"
	DefaultConfigId = "default"
)

func main() {
	var params app.RenderParams

	// init app folders
	app.MakeDir("out")

	if len(os.Args) == 2 && os.Args[1] == "default-image" {
		params = app.ParseFileParams(ExampleConfig, "smooth-boundaries")
	} else if len(os.Args) == 2 && os.Args[1] == "default-video" {
		params = app.ParseFileParams(ExampleConfig, "stripy-video")
	} else if len(os.Args) == 2 {
		// use only parameters defined in config file
		params = app.ParseFileParams(Config, os.Args[1])
	} else if len(os.Args) > 3 {
		// use cli params and overwrite default params
		params = parseCliParams()
	} else {
		// only for development mode
		params = app.ParseFileParams(ExampleConfig, DefaultConfigId)
	}
	app.Render(params)
	fmt.Println("\n --> DONE !")
}

func parseCliParams() app.RenderParams {
	params := app.ParseFileParams("example.config.json", DefaultConfigId)
	params.Resolution = paramToFloat(os.Args[1])
	params.Image.C = complex(paramToFloat(os.Args[2]), paramToFloat(os.Args[3]))
	return params
}

func paramToFloat(input string) float32 {
	n, err := strconv.ParseFloat(input, 32)
	if err != nil {
		panic(fmt.Sprintf("Parameter %s is not a number", input))
	}
	return float32(n)
}
