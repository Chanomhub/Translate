package main

import (
	"fmt"
	"os"
	"path/filepath"

	gt "github.com/bas24/googletranslatefree"
	"github.com/minio/simdjson-go"
)

func main() {
	// Example input JSON
	const inputJSON = `
	{
	  "greeting": "Hello, World!",
	  "description": "This is a sample JSON",
	  "untranslated": "This should remain the same"
	}`

	// Paths for input and output data
	inputFilePath := "input_data.json"  // You'll need to provide the actual file
	outputDir := "translated_output"

	// Translation configuration
	translateKeys := []string{"greeting"}        // Keys to translate
	sourceLang := "en"                          // Source language
	targetLang := "es"                          // Target language

	// -----------------------------------
	// 1. Load and Parse JSON using simdjson-go
	pj, err := simdjson.Parse([]byte(inputJSON), nil)
	if err != nil {
		fmt.Println("JSON parsing error:", err)
		return
	}

	// 2. Iterate and Translate
	iter := pj.Iter()
	for {
		typ := iter.Advance()
		if typ == simdjson.TypeNone {
			break // End of the object
		}

		key, _ := iter.ObjectElementKey()
		shouldTranslate := contains(translateKeys, string(key))

		if shouldTranslate {
			if iter.Type() == simdjson.TypeString {
				originalValue, _ := iter.StringBytes()
				translated, err := gt.Translate(string(originalValue), sourceLang, targetLang)
				if err != nil {
					fmt.Println("Translation error:", err)
				} else {
					iter.ReplaceString(translated)
				}
			} else {
				fmt.Println("Cannot translate non-string value for key:", key)
			}
		}
	}

	// 3. Create output directory if needed
	_ = os.MkdirAll(outputDir, os.ModePerm)

	// 4. Serialize modified JSON
	outputData, _ := pj.MarshalJSON()
	outputFileName := filepath.Base(inputFilePath) // Same name as input
	outputFilePath := filepath.Join(outputDir, outputFileName)

	err = os.WriteFile(outputFilePath, outputData, 0644)
	if err != nil {
		fmt.Println("Error writing output file:", err)
		return
	}

	fmt.Println("Translated JSON saved to:", outputFilePath)
}

// Helper function to check if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
