package app

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
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

type ColorParams struct {
	ColorSpace string
	C1         string
	C2         string
	C3         string
}

type RenderParams struct {
	Id            string
	Resolution    float32
	RenderMode    string
	Encoding      string
	Filename      string
	MaxIterations int
	MaxThreshold  int
	Color         ColorParams
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
		if renderMode == "-t" || renderMode == "T" {
			renderMode = "THRESHOLD"
		} else if renderMode == "-i" || renderMode == "I" {
			renderMode = "ITERATION"
		}
		maxIterations := obj["maxIterations"]
		if maxIterations == nil {
			maxIterations = 20.0
			fmt.Println("maxThreshold not defined (default 20)")
		}
		maxThreshold := obj["maxThreshold"]
		if maxThreshold == nil {
			maxThreshold = 20.0
			fmt.Println("maxThreshold not defined (default 20)")
		}
		encoding := obj["encoding"]
		if encoding == nil {
			panic("Property 'encoding' not defined")
		}
		color := obj["color"]
		var colorSpace string

		var C1 string
		var C2 string
		var C3 string
		if color == nil {
			fmt.Println("color not defined")
		} else {
			c := strings.ReplaceAll(color.(string), " ", "")
			// replace last parenthesis with space
			c = replaceAtIndex(c, ' ', len(c)-1)
			// remove space at the end
			c = strings.ReplaceAll(c, " ", "")
			colorSpace = c[0:3]
			c = strings.Replace(c, colorSpace, "", 1)
			c = strings.Replace(c, "(", "", 1)
			params := strings.Split(c, ",")
			if len(params) < 3 {
				panic(fmt.Sprintf("Invalid color params %s", params))
			}
			C1 = params[0]
			C2 = params[1]
			C3 = params[2]
		}
		filename := obj["filename"]
		if filename == nil {
			filename = currId
		}
		if outputFile != "" {
			filename = outputFile
		}
		renderParams := RenderParams{
			Id:            currId.(string),
			Resolution:    float32(resolution.(float64)),
			RenderMode:    renderMode.(string),
			Encoding:      encoding.(string),
			Filename:      filename.(string),
			MaxIterations: int(maxIterations.(float64)),
			MaxThreshold:  int(maxThreshold.(float64)),
		}
		if color != nil {
			renderParams.Color = ColorParams{
				ColorSpace: colorSpace,
				C1:         C1,
				C2:         C2,
				C3:         C3,
			}
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
			renderParams.Filename += encodeParams(renderParams, true)
		} else if val, ok := obj["static"]; ok {
			renderParams.Image = parseAbstractParams(parseJsonObj(val))
			renderParams.Filename += encodeParams(renderParams, false)
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

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
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

func encodeParams(params RenderParams, isVideo bool) string {
	if isVideo {
		return fmt.Sprintf("_%s_%s_%s", params.RenderMode, encodeColor(params.Color), encodeComplex(params.Video.Start.C))
	} else {
		return fmt.Sprintf("_%s_%s_%s", params.RenderMode, encodeColor(params.Color), encodeComplex(params.Image.C))
	}
}

func encodeColor(color ColorParams) string {
	return fmt.Sprintf("(%s_%s_%s_%s)", color.ColorSpace, color.C1, color.C2, color.C3)
}

func encodeComplex(c complex64) string {
	re := real(c)
	im := imag(c)
	return fmt.Sprintf("(%f%f)", re, im)
}

func PrintParams(params RenderParams) {
	fmt.Printf("Resolution: %f\n", params.Resolution)
	fmt.Printf("Render mode: %s\n", params.RenderMode)
	fmt.Printf("Encoding: %s\n", params.Encoding)
	fmt.Printf("Filename: %s\n", params.Filename)
	fmt.Printf("Max threshold: %d\n", params.MaxThreshold)
	fmt.Printf("Max iterations: %d\n", params.MaxIterations)
	fmt.Printf("Color space: %s\n", params.Color.ColorSpace)
	fmt.Printf("First color param: %s\n", params.Color.C1)
	fmt.Printf("Second color param: %s\n", params.Color.C2)
	fmt.Printf("Third color param: %s\n", params.Color.C3)
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
