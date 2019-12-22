package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func toFloat(input string) float32 {
	n, err := strconv.ParseFloat(input, 32)
	if err != nil {
		panic(fmt.Sprintf("Parameter %s is not a number", input))
	}
	return float32(n)
}

func printPrams(args []string) {
	fmt.Printf("Image size: %s \n", args[1])
	fmt.Printf("Hyperparam C: %s + %si \n", args[2], args[3])
	fmt.Printf("Generation mode: %s \n", args[4])
	fmt.Printf("Output file: %s \n", args[5])
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func write(filename string, data string) {
	output := []byte(data)
	err := ioutil.WriteFile("../outputs/"+filename, output, 0644)
	if err != nil {
		panic(err)
	}
}
