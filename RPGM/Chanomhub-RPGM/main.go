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
	Frames         [][][]float64 `json:"frames"`
	Name           string     `json:"name"`
	Position       int        `json:"position"`
	Timings        []Timings  `json:"timings"`
}

type Timings struct {
	FlashColor    []int `json:"flashColor"`
	FlashDuration int   `json:"flashDuration"`
	FlashScope    int   `json:"flashScope"`
	Frame         int   `json:"frame"`
	Se            Se    `json:"se"`
}

type Se struct {
	Name   string `json:"name"`
	Pan    int    `json:"pan"`
	Pitch  int    `json:"pitch"`
	Volume int    `json:"volume"`
}

// Define a map to associate file names with struct types
var structMap = map[string]interface{}{
	"Animations.json": EnemyData{},
	// Add more mappings here if needed
}

func main() {
	// Command-line flags
	filePath := flag.String("file", "", "Path to the input JSON file")
	targetLang := flag.String("lang", "th", "Target language code (e.g., 'th' for Thai)")
	flag.Parse()

	// Error handling for flags
	if *filePath == "" {
		fmt.Println("Error: Please provide a file path using the -file flag")
		return
	}

	// Load the JSON file
	enemies, err := loadJSONFile(*filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Translate and update names if applicable
	err = translateAndUpdateNames(enemies, *targetLang)
	if err != nil {
		fmt.Println("Translation error:", err)
		return
	}

	// Save the modified JSON
	outputFilePath := filepath.Join("output", filepath.Base(*filePath))
	err = saveJSON(outputFilePath, enemies)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}

	fmt.Println("Translation complete! Results saved to:", outputFilePath)
}

// Function to load JSON file
func loadJSONFile(filePath string) ([]EnemyData, error) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Determine the struct type based on the file name
	structType, ok := structMap[filepath.Base(filePath)]
	if !ok {
		return nil, fmt.Errorf("unsupported file type")
	}

	// Unmarshal (decode) the JSON into the appropriate struct type
	err = json.Unmarshal(jsonData, &structType)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	if enemies, ok := structType.([]EnemyData); ok {
		return enemies, nil
	}

	return nil, fmt.Errorf("unexpected struct type")
}

// Function to translate and update names
func translateAndUpdateNames(enemies []EnemyData, targetLang string) error {
	for i := range enemies {
		if enemies[i].Name != "" {
			translatedText, err := gt.Translate(enemies[i].Name, "auto", targetLang)
			if err != nil {
				return fmt.Errorf("translation error for enemy %s: %w", enemies[i].Name, err)
			}
			enemies[i].Name = translatedText
		}
	}
	return nil
}

// Function to save JSON to file
func saveJSON(filePath string, data interface{}) error {
	updatedJSON, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	err = os.WriteFile(filePath, updatedJSON, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}
