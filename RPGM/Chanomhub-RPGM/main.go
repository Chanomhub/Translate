package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bas24/googletranslatefree"
	"github.com/minio/simdjson-go"
)

// JSONData represents the structure of your JSON data
type JSONData struct {
	// Define fields of your JSON data structure
	// For example:
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	const inputFilePath = "input.json"
	const outputFolder = "translated_json"

	// Read JSON data from file
	jsonData, err := readJSONFile(inputFilePath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// Translate JSON data
	translatedData := translateJSON(jsonData, "en", "es")

	// Export translated JSON data
	outputFileName := strings.TrimSuffix(inputFilePath, filepath.Ext(inputFilePath)) + "_translated.json"
	outputFilePath := filepath.Join(outputFolder, outputFileName)
	err = exportJSON(translatedData, outputFilePath)
	if err != nil {
		fmt.Println("Error exporting translated JSON data:", err)
		return
	}

	fmt.Println("Translated JSON data exported to:", outputFilePath)
}

// readJSONFile reads JSON data from a file
func readJSONFile(filePath string) (JSONData, error) {
	var data JSONData

	file, err := os.Open(filePath)
	if err != nil {
		return data, err
	}
	defer file.Close()

	// Parse JSON data using simdjson-go
	parsedData, err := simdjson.Parse(file, nil)
	if err != nil {
		return data, err
	}

	// Map parsed JSON data to struct
	err = parsedData.Unmarshal(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// translateJSON translates specified fields of JSON data
func translateJSON(data JSONData, sourceLang, targetLang string) JSONData {
	// Translate specified fields (e.g., Title, Description)
	translatedTitle, _ := googletranslatefree.Translate(data.Title, sourceLang, targetLang)
	translatedDescription, _ := googletranslatefree.Translate(data.Description, sourceLang, targetLang)

	// Return translated JSON data
	return JSONData{
		Title:       translatedTitle,
		Description: translatedDescription,
		// Add other fields as needed
	}
}

// exportJSON exports JSON data to a file
func exportJSON(data JSONData, filePath string) error {
	// Create output folder if not exists
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}

	// Marshal JSON data
	jsonBytes, err := simdjson.Marshal(data)
	if err != nil {
		return err
	}

	// Write JSON data to file
	err = ioutil.WriteFile(filePath, jsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
