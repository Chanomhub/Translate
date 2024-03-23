package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	gt "github.com/bas24/googletranslatefree"
)

// Define a function to handle deep translation within JSON structures
func translateJSON(data interface{}, sourceLang, targetLang string) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}: 
		for key, value := range v {
			translatedValue, err := translateJSON(value, sourceLang, targetLang)
			if err != nil {
				return nil, err
			}
			v[key] = translatedValue
		}
		return v, nil

	case []interface{}: 
		for i, value := range v {
			if value != nil { // Only attempt translation if not nil
				translatedValue, err := translateJSON(value, sourceLang, targetLang)
				if err != nil {
					return nil, err
				}
				v[i] = translatedValue
			}
		}
		return v, nil

	case string: 
		translated, err := gt.Translate(v, sourceLang, targetLang)
		if err != nil {
			return nil, fmt.Errorf("error translating: %w", err)
		}
		return translated, nil

	default: 
		return data, nil
	}
}

func main() {
	// Define flags for input file, output file, source/target languages
	inputFileName := flag.String("input", "", "Input JSON file name")
	outputFileName := flag.String("output", "", "Output JSON file name")
	sourceLang := flag.String("source", "auto", "Source language (use 'auto' for detection)")
	targetLang := flag.String("target", "en", "Target language")
	flag.Parse()

	// Input/Output Logic
	if *inputFileName == "" {
		fmt.Println("Please provide an input file name using -input flag")
		return
	}

	jsonFile, err := os.Open(*inputFileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Unmarshal the JSON data
	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Perform translation
	translatedData, err := translateJSON(data, *sourceLang, *targetLang)
	if err != nil {
		fmt.Println("Error during translation:", err)
		return
	}

	// Marshal the translated data back into JSON
	updatedJSON, err := json.MarshalIndent(translatedData, "", " ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Get the output file or use input file as default
	if *outputFileName == "" {
		outputFileName = inputFileName
	}

	// Write the translated JSON back to the file
	err = ioutil.WriteFile(*outputFileName, updatedJSON, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Updated JSON saved to file:", *outputFileName)
}
