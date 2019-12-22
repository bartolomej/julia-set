package app

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type VideoParams struct {
	Start AbstractParams
	End   AbstractParams
}

type AbstractParams struct {
	C        complex64
	CenterX  float32
	CenterY  float32
	AxisSpan float32
}

type RenderParams struct {
	Resolution float32
	RenderMode string
	Encoding   string
	Filename   string
	Image      AbstractParams
	Video      VideoParams
}

func ParseFileParams(id string) RenderParams {
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
		encoding := obj["encoding"]
		if encoding == nil {
			panic("Property 'encoding' not defined")
		}
		filename := obj["filename"]
		if filename == nil {
			filename = fmt.Sprintf("out_%s", currId)
		}
		renderParams := RenderParams{
			Resolution: float32(resolution.(float64)),
			RenderMode: renderMode.(string),
			Encoding:   encoding.(string),
			Filename:   filename.(string),
		}
		start := obj["start"]
		end := obj["end"]
		if start != nil && end != nil {
			renderParams.Video.Start = parseAbstractParams(parseJsonObj(start))
			renderParams.Video.End = parseAbstractParams(parseJsonObj(end))
		} else if val, ok := obj["static"]; ok {
			renderParams.Image = parseAbstractParams(parseJsonObj(val))
		} else {
			panic("No keys matching 'start', 'end' or 'static' params")
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
