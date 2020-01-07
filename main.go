package main

import (
	"fmt"
	"github.com/bartolomej/complex-set-art/app"
	"os"
)

func main() {
	var params app.RenderParams
	// init output folder if doesn't exist
	app.MakeDir("out")
	if len(os.Args) == 2 && os.Args[1] == "default-image" {
		params = getDefaultImageParams()
	} else if len(os.Args) == 2 && os.Args[1] == "default-video" {
		panic("Video rendering not implemented yet")
	} else if len(os.Args) == 2 {
		params = app.ParseFileParams(os.Args[1], "")
	} else if len(os.Args) == 3 {
		params = app.ParseFileParams(os.Args[1], os.Args[2])
	} else if len(os.Args) > 4 {
		params = app.ParseCliParams()
	} else {
		params = app.ParseFileParams("zoomed-in", "")
	}
	// print currently used params
	app.PrintParams(params)
	if (params.Video != app.VideoParams{}) {
		panic("Video rendering not implemented yet")
	} else if (params.Image != app.AbstractParams{}) {
		app.RenderImage(params)
	}
	fmt.Println("DONE !")
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
