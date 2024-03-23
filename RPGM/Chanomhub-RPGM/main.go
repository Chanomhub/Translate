package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	gt "github.com/bas24/googletranslatefree"
)

var (
	// Define flags for input and output file names
	inputFileName  = flag.String("input", "", "Input JSON file name")
	startLanguage  = flag.String("start", "auto", "Starting language (default: auto)")
	targetLanguage = flag.String("target", "en", "Target language (default: en)")
)

// Define a function to handle deep translation within a JSON structure
func translateJSON(data interface{}, startLang, targetLang string) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}: // If it's a map (object)
		for key, value := range v {
			if key == "name" {
				translated, err := gt.Translate(value.(string), startLang, targetLang)
				if err != nil {
					return nil, fmt.Errorf("error translating 'name': %w", err)
				}
				v[key] = translated
			} else {
				// Recursive call for nested structures
				translatedValue, err := translateJSON(value, startLang, targetLang)
				if err != nil {
					return nil, err
				}
				v[key] = translatedValue
			}
		}
		return v, nil
	case []interface{}: // If it's an array
		for i, value := range v {
			translatedValue, err := translateJSON(value, startLang, targetLang)
			if err != nil {
				return nil, err
			}
			v[i] = translatedValue
		}
		return v, nil
	default: // Base case for other data types
		return data, nil
	}
}

func main() {
	// Parse command line flags
	flag.Parse()

	if *inputFileName == "" {
		fmt.Println("Please provide an input file name using -input flag")
		return
	}

	// Read the JSON from the file
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

	// Unmarshal the JSON
	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Translate
	translatedData, err := translateJSON(data, *startLanguage, *targetLanguage)
	if err != nil {
		fmt.Println("Error during translation:", err)
		return
	}

	// Re-encode the updated JSON
	updatedJSON, err := json.MarshalIndent(translatedData, "", " ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Get the output file name from the input file name
	outputFileName := *inputFileName

	// Write back to the file (overwriting it)
	err = ioutil.WriteFile(outputFileName, updatedJSON, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Updated JSON saved to file:", outputFileName)
}
