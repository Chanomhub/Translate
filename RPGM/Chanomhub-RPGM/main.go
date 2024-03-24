package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	gt "github.com/bas24/googletranslatefree"
	"github.com/minio/simdjson-go"
)

func main() {
	jsonFile, err := os.Open("your_rpg_maker_file.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer jsonFile.Close()

	byteData, _ := ioutil.ReadAll(jsonFile)

	// Fast JSON Parsing
	parsed, err := simdjson.Parse(byteData, nil)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Find text elements to translate (adjust selectors as needed)
	elementsToTranslate, _ := parsed.Find("dialog", "text")

	for _, element := range elementsToTranslate.Iter() {
		originalText := element.String()

		translatedText, err := gt.Translate(originalText, "auto", "es") // Translate to Spanish
		if err != nil {
			fmt.Println("Translation error:", err)
		} else {
			element.SetString(translatedText) // Replace text in parsed JSON
		}
	}

	// Overwrite the original file (or save a new one)
	updatedJSON, _ := parsed.MarshalJSON()
	err = ioutil.WriteFile("your_rpg_maker_file_translated.json", updatedJSON, 0644)
	if err != nil {
		fmt.Println("Error writing translated file:", err)
	}
}
