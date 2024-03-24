package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	gt "github.com/bas24/googletranslatefree"
)

type EnemyData struct {
	ID             int        `json:"id"`
	Animation1Hue  int        `json:"animation1Hue"`
	Animation1Name string     `json:"animation1Name"`
	Animation2Hue  int        `json:"animation2Hue"`
	Animation2Name string     `json:"animation2Name"`
	Frames         [][][]int  `json:"frames"`
	Name           string     `json:"name"`
	Position       int        `json:"position"`
	Timings        []Timings  `json:"timings"`
}

type Timings struct {
	FlashColor   []int      `json:"flashColor"`
	FlashDuration int       `json:"flashDuration"`
	FlashScope    int        `json:"flashScope"`
	Frame         int        `json:"frame"`
	Se            Se         `json:"se"`
}

type Se struct {
	Name   string `json:"name"`
	Pan    int    `json:"pan"`
	Pitch  int    `json:"pitch"`
	Volume int    `json:"volume"`
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
	err = os.MkdirAll(outputDir, 0644)
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
