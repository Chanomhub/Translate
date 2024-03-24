package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	gt "github.com/bas24/googletranslatefree"
	"github.com/goccy/go-json"
	"github.com/minio/simdjson-go"
)

// config holds the configurable parameters of the translation process
type config struct {
	inputFilePath     string
	outputDirPath     string
	translationPoints []string
	targetLanguage    string
}

func main() {
	// Parse command-line flags
	cfg := parseFlags()

	// Input validation
	if err := validateConfig(cfg); err != nil {
		fmt.Println(err)
		return
	}

	// Read and parse the input JSON file
	inputBytes, err := ioutil.ReadFile(cfg.inputFilePath)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}
	parsed, err := simdjson.Parse(inputBytes, nil)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Perform the translations
	if err := translateJSON(parsed, cfg.translationPoints, cfg.targetLanguage); err != nil {
		fmt.Println("Translation error:", err)
		return
	}

	// Serialize the translated JSON
	outputBytes, err := json.Marshal(parsed.Iter())
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Create the output file
	if err := writeOutputFile(cfg.outputDirPath, cfg.inputFilePath, outputBytes); err != nil {
		fmt.Println("Error writing output file:", err)
		return
	}

	fmt.Println("Translation complete!")
}

// parseFlags parses command-line flags and returns a config struct
func parseFlags() *config {
	var cfg config
	flag.StringVar(&cfg.inputFilePath, "input", "", "Path to the input JSON file")
	flag.StringVar(&cfg.outputDirPath, "output", "translations", "Directory to store translated JSON files")
	flag.StringSliceVar(&cfg.translationPoints, "keys", []string{"text"}, "Comma-separated list of keys within the JSON to translate")
	flag.StringVar(&cfg.targetLanguage, "lang", "es", "Target language code (e.g., es, fr, de)")
	flag.Parse()
	return &cfg
}

// validateConfig checks the provided configuration for validity.
func validateConfig(cfg *config) error {
	if cfg.inputFilePath == "" {
		return fmt.Errorf("input file path is required")
	}
	// ... add other validation checks as needed 
	return nil
}

// translateJSON performs translations on the parsed simdjson object
func translateJSON(parsed *simdjson.ParsedJson, keys []string, targetLanguage string) error {
	for _, key := range keys {
		valueIter := parsed.FindKey(key)
		if valueIter != nil {
			value, err := valueIter.String()
			if err != nil {
				return fmt.Errorf("error fetching translatable string: %w", err)
			}

			translatedText, err := gt.Translate(value, "auto", targetLanguage)
if err != nil {
	// Handle translation error here (e.g., log the error, skip translation)
	fmt.Printf("Error translating key '%s': %v\n", key, err)
	// You can choose to skip translation for this key or continue translating others
	continue
}



			// Update JSON in-place using simdjson-go (implementation omitted for brevity)
			// ...
		}
	}
	return nil
}

// writeOutputFile creates the output file with the translated JSON content
func writeOutputFile(outputDir, inputFilename string, data []byte) error {
	outputFilename := filepath.Base(inputFilename)
	outputFilePath := filepath.Join(outputDir, outputFilename)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	if err := ioutil.WriteFile(outputFilePath, data, 0644); err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}
	return nil
}
