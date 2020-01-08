package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
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
	Exponent complex64
	OriginX  float32
	OriginY  float32
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
	ReturnMode    ReturnMode
	Encoding      string
	Filename      string
	Folder        string
	MaxIterations int
	MaxThreshold  int
	Color         ColorParams
	Image         AbstractParams
	Video         VideoParams
}

func ParseFileParams(configId string, outputFile string) RenderParams {
	rawJsonConfig, err := ReadFile("renders.json")
	if err != nil {
		panic("Error opening renders.json config file")
	}

	var jsonConfig []interface{}
	_ = json.Unmarshal([]byte(rawJsonConfig), &jsonConfig)
	for i := 0; i < len(jsonConfig); i++ {
		config := ParseJsonObject(jsonConfig[i])
		var configErrors []error

		id, idErr := getConfigProp(config, "id")
		resolution, resolutionErr := getConfigProp(config, "resolution")
		filename, _ := getConfigProp(config, "filename")
		returnMode, renderModeErr := getConfigProp(config, "returnMode")
		maxIterations, _ := getConfigProp(config, "maxIterations")
		maxThreshold, _ := getConfigProp(config, "maxThreshold")
		encoding, encodingErr := getConfigProp(config, "encoding")
		color, _ := getConfigProp(config, "color")
		colorParams, colorParamsErr := parseColorParams(color)

		// if id doesn't match skip parsing and try next config
		if id != configId {
			continue
		}

		if idErr != nil {
			configErrors = append(configErrors, idErr)
		}
		if renderModeErr != nil {
			configErrors = append(configErrors, renderModeErr)
		}
		if encodingErr != nil {
			configErrors = append(configErrors, encodingErr)
		}
		if resolutionErr != nil {
			configErrors = append(configErrors, resolutionErr)
		}
		if colorParamsErr != nil {
			configErrors = append(configErrors, colorParamsErr)
		}

		if filename == nil {
			filename = id
		}
		if outputFile != "" {
			filename = outputFile
		}
		if maxIterations == nil {
			maxIterations = 20.0
		}
		if maxThreshold == nil {
			maxThreshold = 20.0
		}

		renderParams := RenderParams{
			Id:            id.(string),
			Resolution:    float32(resolution.(float64)),
			ReturnMode:    ParseReturnMode(returnMode.(string)),
			Encoding:      encoding.(string),
			Filename:      filename.(string),
			MaxIterations: int(maxIterations.(float64)),
			MaxThreshold:  int(maxThreshold.(float64)),
			Color:         colorParams,
		}

		static, _ := getConfigProp(config, "static")
		if static != nil {
			renderParams.Image = parseAbstractParams(ParseJsonObject(static))
			renderParams.Filename += encodeParams(renderParams, false)
			handleErrors(configErrors)
			printParams(renderParams)
			return renderParams
		}

		start, startErr := getConfigProp(config, "start")
		end, endErr := getConfigProp(config, "end")
		fps, fpsErr := getConfigProp(config, "fps")
		duration, durationErr := getConfigProp(config, "duration")

		if startErr != nil {
			configErrors = append(configErrors, startErr)
		}
		if endErr != nil {
			configErrors = append(configErrors, endErr)
		}
		if fpsErr != nil {
			configErrors = append(configErrors, fpsErr)
		}
		if durationErr != nil {
			configErrors = append(configErrors, durationErr)
		}

		renderParams.Video = VideoParams{
			Fps:      int(fps.(float64)),
			Duration: duration.(float64),
			Start:    parseAbstractParams(ParseJsonObject(start)),
			End:      parseAbstractParams(ParseJsonObject(end)),
		}

		renderParams.Filename += encodeParams(renderParams, true)
		handleErrors(configErrors)
		printParams(renderParams)
		return renderParams
	}
	panic(fmt.Sprintf("Config with id: %s not found", configId))
}

func handleErrors(errors []error) {
	if len(errors) == 0 {
		return
	}
	fmt.Println("CONFIGURATION ERRORS:")
	for i := 0; i < len(errors); i++ {
		fmt.Println(errors[i].Error())
	}
	os.Exit(1)
}

func ParseReturnMode(returnMode string) ReturnMode {
	if returnMode == "-d" || returnMode == "D" || returnMode == "DISTANCE" {
		return DISTANCE
	} else if returnMode == "-i" || returnMode == "I" || returnMode == "ITERATION" {
		return ITERATION
	}
	panic("Invalid return mode " + returnMode)
}

func parseColorParams(color interface{}) (ColorParams, error) {
	if color == nil {
		return ColorParams{}, nil
	} else {
		c := strings.ReplaceAll(color.(string), " ", "")
		// replace last parenthesis with space
		c = replaceAtIndex(c, ' ', len(c)-1)
		// remove space at the end
		c = strings.ReplaceAll(c, " ", "")
		colorSpace := c[0:3]
		c = strings.Replace(c, colorSpace, "", 1)
		c = strings.Replace(c, "(", "", 1)
		params := strings.Split(c, ",")
		if len(params) < 3 {
			return ColorParams{}, errors.New(fmt.Sprintf("Invalid color params %s", params))
		} else {
			colorParams := ColorParams{
				ColorSpace: colorSpace,
				C1:         params[0],
				C2:         params[1],
				C3:         params[2],
			}
			return colorParams, nil
		}
	}
}

func parseAbstractParams(obj map[string]interface{}) AbstractParams {
	var absParamsErrors []error

	originX, originXErr := getConfigProp(obj, "originX")
	originY, centerYErr := getConfigProp(obj, "originY")
	axisSpan, axisSpanErr := getConfigProp(obj, "axisSpan")
	realC, realCErr := getConfigProp(obj, "realC")
	imagC, imagCErr := getConfigProp(obj, "imagC")
	realExp, realExpErr := getConfigProp(obj, "realExp")
	imagExp, imagExpErr := getConfigProp(obj, "imagExp")

	if originXErr != nil {
		absParamsErrors = append(absParamsErrors, originXErr)
	}
	if centerYErr != nil {
		absParamsErrors = append(absParamsErrors, centerYErr)
	}
	if axisSpanErr != nil {
		absParamsErrors = append(absParamsErrors, axisSpanErr)
	}
	if realCErr != nil {
		absParamsErrors = append(absParamsErrors, realCErr)
	}
	if imagCErr != nil {
		absParamsErrors = append(absParamsErrors, imagCErr)
	}
	if realExpErr != nil {
		absParamsErrors = append(absParamsErrors, realExpErr)
	}
	if imagExpErr != nil {
		absParamsErrors = append(absParamsErrors, imagExpErr)
	}

	handleErrors(absParamsErrors)

	c := complex64(complex(realC.(float64), imagC.(float64)))
	exp := complex64(complex(realExp.(float64), imagExp.(float64)))

	return AbstractParams{
		OriginX:  float32(originX.(float64)),
		OriginY:  float32(originY.(float64)),
		AxisSpan: float32(axisSpan.(float64)),
		Exponent: exp,
		C:        c,
	}
}

func getConfigProp(config map[string]interface{}, prop string) (interface{}, error) {
	resolution := config[prop]
	if resolution == nil {
		return nil, errors.New(fmt.Sprintf("Property '%s' not defined", prop))
	} else {
		return resolution, nil
	}
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func encodeParams(params RenderParams, isVideo bool) string {
	if isVideo {
		return fmt.Sprintf("_%s_%s", params.ReturnMode, encodeComplex(params.Video.Start.C))
	} else {
		return fmt.Sprintf("_%s_%s", params.ReturnMode, encodeComplex(params.Image.C))
	}
}

func encodeComplex(c complex64) string {
	re := real(c)
	im := imag(c)
	return fmt.Sprintf("(%f%f)", re, im)
}

func printParams(params RenderParams) {
	fmt.Println()
	fmt.Println("PARAMETERS: ")
	fmt.Printf("Resolution: %f\n", params.Resolution)
	fmt.Printf("Render mode: %s\n", params.ReturnMode)
	fmt.Printf("Encoding: %s\n", params.Encoding)
	fmt.Printf("Filename: %s\n", params.Filename)
	fmt.Printf("Max threshold: %d\n", params.MaxThreshold)
	fmt.Printf("Max iterations: %d\n", params.MaxIterations)
	fmt.Printf("Color: %s\n", encodeColorForPrint(params.Color))
	if (params.Image != AbstractParams{}) {
		printAbstractParams(params.Image)
	} else if (params.Video != VideoParams{}) {
		fmt.Printf("Duration %f\n", params.Video.Duration)
		fmt.Printf("Fps: %d\n", params.Video.Fps)
		fmt.Println("START: ")
		printAbstractParams(params.Video.Start)
		fmt.Println("END: ")
		printAbstractParams(params.Video.End)
	}
	fmt.Println()
}

func printAbstractParams(params AbstractParams) {
	fmt.Printf("C: %f\n", params.C)
	fmt.Printf("Exponent: %f\n", params.Exponent)
	fmt.Printf("OriginX: %f\n", params.OriginX)
	fmt.Printf("OriginY: %f\n", params.OriginY)
	fmt.Printf("Axis span: %f\n", params.AxisSpan)
}

func encodeColorForPrint(color ColorParams) string {
	return fmt.Sprintf("%s(%s, %s, %s)", color.ColorSpace, color.C1, color.C2, color.C3)
}
