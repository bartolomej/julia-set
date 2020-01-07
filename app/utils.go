package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func ParamToFloat(input string) float32 {
	n, err := strconv.ParseFloat(input, 32)
	if err != nil {
		panic(fmt.Sprintf("Parameter %s is not a number", input))
	}
	return float32(n)
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func ReadFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	return string(b), nil
}

func WriteFile(filename string, data string) {
	output := []byte(data)
	err := ioutil.WriteFile("../outputs/"+filename, output, 0644)
	if err != nil {
		panic(err)
	}
}

func MakeDir(path string) {
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		if !strings.Contains(err.Error(), "file exists") {
			panic(err)
		}
	}
}
