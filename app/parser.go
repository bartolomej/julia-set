package app

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type VideoParams struct {
	Fps      int
	Duration float64
	Start    AbstractParams
	End      AbstractParams
}

type AbstractParams struct {
	C        complex64
	CenterX  float32
	CenterY  float32
	AxisSpan float32
}

type RenderParams struct {
	Resolution    float32
	RenderMode    string
	Encoding      string
	Filename      string
	MaxIterations int
	MaxThreshold  int
	Image         AbstractParams
	Video         VideoParams
}

func ParseFileParams(id string, outputFile string) RenderParams {
	str, err := ReadFile("renders.json")
	if err != nil {
		panic("Error opening renders.json config file")
	}
	var res []interface{}
	_ = json.Unmarshal([]byte(str), &res)
	for i := 0; i < len(res); i++ {
		obj := parseJsonObj(res[i])
		currId := obj["id"]
		if currId == nil {
			panic("Property 'id' not defined")
		}
		resolution := obj["resolution"]
		if resolution == nil {
			panic("Property 'resolution' not defined")
		}
		renderMode := obj["renderMode"]
		if renderMode == nil {
			panic("Property 'renderMode' not defined")
		}
		maxIterations := obj["maxIterations"]
		if maxIterations == nil {
			maxIterations = 20.0
		}
		maxThreshold := obj["maxThreshold"]
		if maxThreshold == nil {
			maxThreshold = 20.0
		}
		encoding := obj["encoding"]
		if encoding == nil {
			panic("Property 'encoding' not defined")
		}
		filename := obj["filename"]
		if filename == nil {
			filename = fmt.Sprintf("out/out_%s", currId)
		} else {
			filename = fmt.Sprintf("out/%s", filename)
		}
		if outputFile != "" {
			filename = fmt.Sprintf("out/%s", outputFile)
		}
		renderParams := RenderParams{
			Resolution:    float32(resolution.(float64)),
			RenderMode:    renderMode.(string),
			Encoding:      encoding.(string),
			Filename:      filename.(string),
			MaxIterations: int(maxIterations.(float64)),
			MaxThreshold:  int(maxThreshold.(float64)),
		}
		start := obj["start"]
		end := obj["end"]
		fps := obj["fps"]
		duration := obj["duration"]
		if start != nil && end != nil && fps != nil && duration != nil {
			renderParams.Video = VideoParams{
				Fps:      int(fps.(float64)),
				Duration: duration.(float64),
				Start:    parseAbstractParams(parseJsonObj(start)),
				End:      parseAbstractParams(parseJsonObj(end)),
			}
			renderParams.Filename += encodeComplex(renderParams.Video.Start.C)
		} else if val, ok := obj["static"]; ok {
			renderParams.Image = parseAbstractParams(parseJsonObj(val))
			renderParams.Filename += encodeComplex(renderParams.Image.C)
		} else {
			panic("Invalid configuration for video/image settings")
		}
		if currId == id {
			return renderParams
		}
	}
	panic(fmt.Sprintf("Config with id: %s not found", id))
}

func ParseCliParams() RenderParams {
	res := ParamToFloat(os.Args[1])
	c := complex(ParamToFloat(os.Args[2]), ParamToFloat(os.Args[3]))
	filename := fmt.Sprintf("r%fi%f_%s", real(c), imag(c), os.Args[4])
	// use static config for cli args
	image := AbstractParams{
		C:        c,
		CenterX:  0,
		CenterY:  0,
		AxisSpan: 2,
	}
	return RenderParams{
		Resolution: res,
		RenderMode: os.Args[4],
		Encoding:   "png",
		Filename:   filename,
		Image:      image,
	}
}

func parseAbstractParams(obj map[string]interface{}) AbstractParams {
	centerX := obj["centerX"].(float64)
	centerY := obj["centerY"].(float64)
	axisSpan := obj["axisSpan"].(float64)
	c := complex64(complex(obj["realC"].(float64), obj["imagC"].(float64)))
	return AbstractParams{
		C:        c,
		CenterX:  float32(centerX),
		CenterY:  float32(centerY),
		AxisSpan: float32(axisSpan),
	}
}

func parseJsonObj(jsonObj interface{}) map[string]interface{} {
	obj := make(map[string]interface{})
	v := reflect.ValueOf(jsonObj)
	if v.Kind() != reflect.Map {
		panic("Json config not of type array")
	}
	for _, key := range v.MapKeys() {
		k := key.Interface().(string)
		obj[k] = v.MapIndex(key).Interface()
	}
	return obj
}

func encodeComplex(c complex64) string {
	re := real(c)
	im := imag(c)
	return fmt.Sprintf("_%f_%f", re, im)
}

func PrintParams(params RenderParams) {
	fmt.Printf("Resolution: %f\n", params.Resolution)
	fmt.Printf("Render mode: %s\n", params.RenderMode)
	fmt.Printf("Encoding: %s\n", params.Encoding)
	fmt.Printf("Filename: %s\n", params.Filename)
	fmt.Printf("Max threshold: %d\n", params.MaxThreshold)
	fmt.Printf("Max iterations: %d\n", params.MaxIterations)
	if (params.Image != AbstractParams{}) {
		printAbstractParams(params.Image)
	} else if (params.Video != VideoParams{}) {
		fmt.Printf("Duration %d\n", params.Video.Duration)
		fmt.Printf("Fps: %d\n", params.Video.Fps)
		fmt.Println("START: ")
		printAbstractParams(params.Video.Start)
		fmt.Println("END: ")
		printAbstractParams(params.Video.End)
	}
}

func printAbstractParams(params AbstractParams) {
	fmt.Printf("C: %f\n", params.C)
	fmt.Printf("CenterX: %f\n", params.CenterX)
	fmt.Printf("CenterY: %f\n", params.CenterY)
	fmt.Printf("Axis span: %f\n", params.AxisSpan)
}
