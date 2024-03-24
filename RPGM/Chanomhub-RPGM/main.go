package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	gt "github.com/bas24/googletranslatefree"
)

// Function to translate the JSON (remains mostly the same)
func translateJSON(data interface{}, targetLanguage string) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			if key == "name" {
				translated, err := gt.Translate(value.(string), "auto", targetLanguage)
				if err != nil {
					return nil, fmt.Errorf("error translating 'name': %w", err)
				}
				v[key] = translated
			} else {
				translatedValue, err := translateJSON(value, targetLanguage)
				if err != nil {
					return nil, err
				}
				v[key] = translatedValue
			}
		}
	case []interface{}:
		for i, value := range v {
			translatedValue, err := translateJSON(value, targetLanguage)
			if err != nil {
				return nil, err
			}
			v[i] = translatedValue
		}
	default:
		// No translation needed for other types
		return data, nil
	}
	return data, nil
}

func main() {
	inputFileName := flag.String("input", "", "Input JSON file name")
	flag.Parse()

	if *inputFileName == "" {
		fmt.Println("Please provide an input file name using -input flag")
		return
	}

	// 1. Open both the input JSON and the configuration file
	jsonFile, err := os.Open(*inputFileName)
	if err != nil {
		fmt.Println("Error opening input JSON file:", err)
		return
	}
	defer jsonFile.Close()

	configureFile, err := os.Open("Configure.json")
	if err != nil {
		fmt.Println("Error opening configuration file:", err)
		return
	}
	defer configureFile.Close()

	// 2. Read data from both JSON files
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading input JSON file:", err)
		return
	}

	configureData, err := ioutil.ReadAll(configureFile)
	if err != nil {
		fmt.Println("Error reading configuration file:", err)
		return
	}

	// 3. Unmarshal both JSONs
	var inputData interface{}
	err = json.Unmarshal(jsonData, &inputData)
	if err != nil {
		fmt.Println("Error decoding input JSON:", err)
		return
	}

	var config map[string]string
	err = json.Unmarshal(configureData, &config)
	if err != nil {
		fmt.Println("Error decoding configuration JSON:", err)
		return
	}

	// 4. Get the target language
	targetLanguage := config["target_language"]

	// 5. Translate
	translatedData, err := translateJSON(inputData, targetLanguage)
	if err != nil {
		fmt.Println("Error during translation:", err)
		return
	}

	// 6. Re-encode translated JSON
	updatedJSON, err := json.MarshalIndent(translatedData, "", " ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// 7. Write updated JSON to file
	err = ioutil.WriteFile(*inputFileName, updatedJSON, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Updated JSON saved to file:", *inputFileName)
}
