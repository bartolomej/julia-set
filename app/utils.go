package app

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

func ParseJsonObject(jsonObj interface{}) map[string]interface{} {
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
