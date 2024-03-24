package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	gt "github.com/bas24/googletranslatefree"
)

// Define a custom struct to mirror the JSON data
type EnemyData struct {
	ID      int `json:"id"`
	Members []struct {
	} `json:"members"`
	Name  string `json:"name"`
	Pages []struct {
	} `json:"pages"`
}

func main() {
	// Command-line flags
	filePath := flag.String("file", "", "Path to the input JSON file")
	targetLang := flag.String("lang", "es", "Target language code (e.g., 'es' for Spanish)")
	flag.Parse()

	// Error handling for flags
	if *filePath == "" {
		fmt.Println("Error: Please provide a file path using the -file flag")
		return
	}

	// Load the JSON file
	jsonData, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Unmarshal (decode) the JSON into an array of EnemyData structs
	var enemies []EnemyData
	err = json.Unmarshal(jsonData, &enemies)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Translate and update names
	for i := range enemies {
		if enemies[i].Name != "" {
			translatedText, err := gt.Translate(enemies[i].Name, "auto", *targetLang)
			if err != nil {
				fmt.Println("Translation error:", err)
			} else {
				enemies[i].Name = translatedText
			}
		}
	}

	// Create a new directory for output files
	outputDir := filepath.Join(".", "output")
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	// Determine output file name
	baseFileName := filepath.Base(*filePath)
	outputFilePath := filepath.Join(outputDir, baseFileName)

	// Marshal (encode) back into JSON
	updatedJSON, err := json.MarshalIndent(enemies, "", " ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Save the modified JSON
	err = os.WriteFile(outputFilePath, updatedJSON, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Translation complete! Results saved to:", outputFilePath)
}
