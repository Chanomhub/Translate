package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	gt "github.com/bas24/googletranslatefree"
)

// Define a function to handle deep translation within a JSON structure
func translateJSON(data interface{}) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}: // If it's a map (object)
		for key, value := range v {
			if key == "name" {
				translated, err := gt.Translate(value.(string), "auto", "en")
				if err != nil {
					return nil, fmt.Errorf("error translating 'name': %w", err)
				}
				v[key] = translated
			} else {
				// Recursive call for nested structures
				translatedValue, err := translateJSON(value)
				if err != nil {
					return nil, err
				}
				v[key] = translatedValue
			}
		}
		return v, nil
	case []interface{}: // If it's an array
		for i, value := range v {
			translatedValue, err := translateJSON(value)
			if err != nil {
				return nil, err
			}
			v[i] = translatedValue
		}
		return v, nil
	default: // Base case for other data types
		return data, nil
	}
}

func main() {
	// 1. Read the JSON from your file
	jsonFile, err := os.Open("Troops.json")
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

	// 5. Write back to the file (overwriting it)
	err = ioutil.WriteFile("Troops.json", updatedJSON, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Updated JSON saved to file.")
}
