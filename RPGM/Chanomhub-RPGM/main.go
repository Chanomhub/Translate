package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	gt "github.com/bas24/googletranslatefree"
	"github.com/minio/simdjson-go"
)

func translateText(text, sourceLang, targetLang string) (string, error) {
	return gt.Translate(text, sourceLang, targetLang)
}

func translateJSON(filePath, sourceLang, targetLang string, fieldsToTranslate []string) error {
	// Read JSON file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Parse JSON
	parsed, err := simdjson.Parse(data, nil)
	if err != nil {
		return err
	}

	// Traverse and translate specified fields
	for _, field := range fieldsToTranslate {
		val, err := parsed.SearchKeyInsensitive([]byte(field))
		if err == nil && val.Type == simdjson.TypeString {
			translated, err := translateText(string(val.IterBytes()), sourceLang, targetLang)
			if err != nil {
				return err
			}
			err = parsed.SetBytes([]byte(translated), []string{field})
			if err != nil {
				return err
			}
		}
	}

	// Marshal the updated JSON
	updatedJSON, err := parsed.MarshalJSON()
	if err != nil {
		return err
	}

	// Write the updated JSON to a new file
	outputPath := filepath.Join("output", filepath.Base(filePath))
	err = ioutil.WriteFile(outputPath, updatedJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	const (
		sourceLang      = "en"
		targetLang      = "es"
		inputFolder     = "input"
		fieldsToTranslate = []string{"title", "description"} // specify fields to translate
	)

	// Create output directory if it doesn't exist
	outputDir := "output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}

	// Get list of JSON files in input folder
	fileList, err := filepath.Glob(filepath.Join(inputFolder, "*.json"))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Translate each JSON file
	for _, file := range fileList {
		err := translateJSON(file, sourceLang, targetLang, fieldsToTranslate)
		if err != nil {
			fmt.Printf("Error translating file %s: %v\n", file, err)
		} else {
			fmt.Printf("Translation successful for file %s\n", file)
		}
	}
}
