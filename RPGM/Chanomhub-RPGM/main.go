package main

import (
	"fmt"
	"io/ioutil"
	"os"

	gt "github.com/bas24/googletranslatefree"
	"github.com/minio/simdjson-go"
)

// Function to perform the translation
func translateText(text, sourceLang, targetLang string) (string, error) {
	return gt.Translate(text, sourceLang, targetLang)
}

func main() {
	// --- Example Usage ---

	// 1. Read original JSON file
	originalData, err := ioutil.ReadFile("input.json")
	if err != nil {
		fmt.Println("Error reading input JSON:", err)
		return
	}

	// 2. Parse the JSON data (using simdjson-go)
	parsed, err := simdjson.Parse(originalData, nil)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// 3. Iterate and translate specified fields
	iter := parsed.Iter()
	for {
		typ := iter.Advance()
		switch typ {
		case simdjson.TypeObject:
			obj, _ := iter.Object(nil)
			objIter := obj.Iter()
		FIELD_LOOP:
			for {
				field := objIter.Advance()
				switch field {
				case simdjson.TypeString("title"): // Specify the field to translate
					str, _ := objIter.StringBytes()
					translated, err := translateText(string(str), "en", "es") // Example translation
					if err != nil {
						fmt.Println("Translation error:", err)
					} else {
						iter.DelField()          // Remove the original field
						iter.AddString(field, translated) // Add the translated field
					}
				default:
					break FIELD_LOOP // Skip other fields
				}
			}
		default:
			// Handle other JSON types if needed
		}
	}

	// 4. Write the modified JSON to a new file
	outputFilename := "output_es.json" // Use language code for clarity
	outputData := parsed.MarshalJSON()
	err = ioutil.WriteFile(outputFilename, outputData, 0644)
	if err != nil {
		fmt.Println("Error writing output JSON:", err)
		return
	}

	fmt.Println("Translation and JSON export complete!")
}
