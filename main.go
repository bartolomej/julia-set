package main

import (
	"fmt"
	"github.com/bartolomej/complex-set-art/app"
	"os"
)

// smooth coloring: https://iquilezles.org/www/articles/mset_smooth/mset_smooth.htm

func main() {
	var params app.RenderParams
	// init output folder if doesn't exist
	app.MakeDir("out")
	if len(os.Args) == 2 {
		params = app.ParseFileParams(os.Args[1])
	} else if len(os.Args) > 4 {
		params = app.ParseCliParams()
	} else {
		params = getDefaultParams()
	}
	// print currently used params
	fmt.Println(params)
	if (params.Video != app.VideoParams{}) {
		panic("Video rendering not yet supported")
	} else if (params.Image != app.AbstractParams{}) {
		generateImage(params)
	}
	fmt.Println("DONE !")
}

func generateImage(renderParams app.RenderParams) {
	setParams := app.SetParams{
		CenterX:    renderParams.Image.CenterX,
		CenterY:    renderParams.Image.CenterY,
		Resolution: renderParams.Resolution,
		AxisSpan:   renderParams.Image.AxisSpan,
		C:          renderParams.Image.C,
	}
	if renderParams.RenderMode == "-i" {
		set := app.CalcByIterations(setParams)
		app.DrawByIteration(set, renderParams)
	} else if renderParams.RenderMode == "-t" {
		set := app.CalcByThreshold(setParams)
		app.DrawByThreshold(set, renderParams)
	} else {
		panic(fmt.Sprintf("Invalid RenderMode %s", renderParams.RenderMode))
	}
}

func generateVideo(params app.RenderParams) {

}

func getDefaultParams() app.RenderParams {
	return app.RenderParams{
		Resolution: 1000,
		RenderMode: "-i",
		Encoding:   "png",
		Filename:   "test-out",
		Image: app.AbstractParams{
			C:        complex(0, 0),
			CenterX:  0,
			CenterY:  0,
			AxisSpan: 2,
		},
	}
}
