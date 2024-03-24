package main

import (
    "encoding/json"
    "fmt"
    "os"
    "flag" // For command-line flags

    gt "github.com/bas24/googletranslatefree"
)

// Define a custom struct to mirror the JSON data
type EnemyData struct {
    ID      int      `json:"id"`
    Members []struct { // ... (other fields if needed)
    } `json:"members"`
    Name    string   `json:"name"`
    Pages   []struct { // ...
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
            translatedText, err := gt.Translate(enemies[i].Name, "en", *targetLang) 
            if err != nil {
                fmt.Println("Translation error:", err)
            } else {
                enemies[i].Name = translatedText
            }
        }
    }

    // Marshal (encode) back into JSON
    updatedJSON, err := json.MarshalIndent(enemies, "", "  ") // Indentation for readability
    if err != nil {
        fmt.Println("Error encoding JSON:", err)
        return
    }

    // Save the modified JSON (you can choose a different output file name)
    err = os.WriteFile("translated_output.json", updatedJSON, 0644)
    if err != nil {
        fmt.Println("Error writing file:", err)
        return
    }

    fmt.Println("Translation complete! Results saved in translated_output.json")
}
