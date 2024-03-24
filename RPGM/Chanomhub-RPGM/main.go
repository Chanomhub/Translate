package main

import (
	"fmt"
	"os"

	"github.com/bas24/googletranslatefree"
	"github.com/minio/simdjson-go"
)

func main() {
	// Example JSON input
	inputFilename := "data.json"
	outputFolder := "translations"

	// Read the JSON data
	data, err := os.ReadFile(inputFilename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Parse the JSON data using simdjson-go
	pj, err := simdjson.Parse(data, nil)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Specify the fields you want to translate
	fieldsToTranslate := []string{"message", "description"}

	// Translation function
	translateField := func(field string, iter *simdjson.Iter) error {
		// Get original text based on the field name
		currentText, err := iter.FindKey(field, nil).String()
		if err != nil {
			return err
		}

		// Specify target languages
		targetLanguages := []string{"es", "fr", "de"} 

		for _, lang := range targetLanguages {
			result, err := gt.Translate(currentText, "auto", lang)
			if err != nil {
				return fmt.Errorf("translation to %s failed: %v", lang, err)
			}

			// Update the JSON object
			iter.SetString(result)
		}
		return nil
	}

	// Iterate and translate
	iter := pj.Iter()
	for {
		typ := iter.Advance()

		switch typ {
		case simdjson.TypeObject:
			obj, err := iter.Object(nil)
			if err != nil {
				fmt.Println("Error getting object:", err)
				continue
			}

			oIter := obj.Iter()
			for {
				field, t := oIter.Advance()
				if t == simdjson.TypeNone {
					break
				}

				if contains(fieldsToTranslate, field) {
					if err := translateField(field, &oIter); err != nil {
						fmt.Println("Error translating field:", err)
					}
				}
			}

		default:
			// Handle other types if needed, likely nothing to translate
		}

		if typ == simdjson.TypeNone {
			break // We're done!
		}
	}

	// Serialize results with simdjson-go
	translatedJSON, err := pj.MarshalJSON() 
	if err != nil {
		fmt.Println("Failed to serialize translated JSON:", err)
		return
	}

	// Create the output folder if it doesn't exist
	os.MkdirAll(outputFolder, 0755)

	// Save translated JSON file
	outputFilename := fmt.Sprintf("%s/%s", outputFolder, inputFilename)
	err = os.WriteFile(outputFilename, translatedJSON, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Translation completed!")
}

// Helper function (case-insensitive)
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
