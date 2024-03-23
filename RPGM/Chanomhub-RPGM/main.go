package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	gt "github.com/bas24/googletranslatefree"
)

// Revised translateJSON for robustness
func translateJSON(data interface{}) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			if key == "name" {
				switch strValue := value.(type) { // Check type before translating
				case string:
					translated, err := gt.Translate(strValue, *fromLanguage, *targetLanguage)
					if err != nil {
						return nil, fmt.Errorf("error translating 'name': %w", err)
					}
					v[key] = translated
				default:
					fmt.Println("Warning: Non-string value found for 'name' field. Skipping translation.")
				}
			} else {
				translatedValue, err := translateJSON(value)
				if err != nil {
					return nil, err
				}
				v[key] = translatedValue
			}
		}
		return v, nil
	// ... (Similar case for []interface{})
	default:
		return data, nil
	}
}

func main() {
	// Define flags for input file name and target language
	inputFileName := flag.String("input", "i", "Input JSON file name")
	targetLanguage := flag.String("target", "t", "Target language code (e.g., 'es' for Spanish, 'fr' for French)")
        fromLanguage := flag.String("from", "f", "Your default game language or 'auto' can be used.")
	flag.Parse()

	if *inputFileName == "" {
		fmt.Println("Please specify the file name for input using '-i' or '-input'.")
		return
	}

	if *targetLanguage == "" {
		fmt.Println("Please specify the target language using '-target' or '-t'.")
		return
	}

	if *fromLanguage == "" {
		fmt.Println("Please default language using '-f ' or '-from'.")
		return
	}

	// 1. Read the JSON from your file
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

	// 2. Unmarshal the JSON
	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// 3. Translate
	translatedData, err := translateJSON(data)
	if err != nil {
		fmt.Println("Error during translation:", err)
		return
	}

	// 4. Re-encode the updated JSON
	updatedJSON, err := json.MarshalIndent(translatedData, "", " ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Get the output file name from the input file name
	outputFileName := *inputFileName

	// 5. Write back to the file (overwriting it)
	err = ioutil.WriteFile(outputFileName, updatedJSON, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Updated JSON saved to file:", outputFileName)
}
