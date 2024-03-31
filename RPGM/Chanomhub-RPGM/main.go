package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"
  "strings"

  "github.com/Conight/go-googletrans"
)

type Event struct {
  ID       int      `json:"id"`
  Name     string   `json:"name"`
  List     []Action `json:"list"`
  SwitchID int      `json:"switchId"`
  Trigger  int      `json:"trigger"`
}

type Action struct {
  Code       int           `json:"code"`
  Indent     int           `json:"indent"`
  Parameters []interface{} `json:"parameters"`
}

type FilePaths struct {
  Files []string `json:"files"`
}

func main() {
  filepaths, err := readFilePaths("file_paths.json")
  if err != nil {
    fmt.Println("Error reading file paths:", err)
    return
  }

  for _, filePath := range filepaths.Files {
    fmt.Println("Processing file:", filePath)

    events, err := readJSONFile(filePath)
    if err != nil {
      fmt.Printf("Error reading file %s: %v\n", filePath, err)
      continue
    }

    for _, event := range events {
      fmt.Println("Original Event Name:", event.Name)
      // Translate event name
      if translatedName, err := translateText(event.Name, "th"); err != nil {
        fmt.Println("Translation error:", err)
      } else {
        event.Name = translatedName
        fmt.Println("Translated Event Name:", event.Name)
      }

      // Process actions
      processActions(event.List)
    }

    // Export translated JSON for each file
    if err := exportTranslatedJSON(events, filePath); err != nil {
      fmt.Printf("Error exporting translated JSON for file %s: %v\n", filePath, err)
      continue
    }

    fmt.Printf("Translation for file %s completed successfully.\n", filePath)
  }
}

func readFilePaths(filename string) (FilePaths, error) {
  var filePaths FilePaths

  file, err := os.Open(filename)
  if err != nil {
    return filePaths, err
  }
  defer file.Close()

  byteValue, err := ioutil.ReadAll(file)
  if err != nil {
    return filePaths, err
  }

  err = json.Unmarshal(byteValue, &filePaths)
  return filePaths, err
}

func readJSONFile(filename string) ([]Event, error) {
  jsonFile, err := os.Open(filename)
  if err != nil {
    return nil, err
  }
  defer jsonFile.Close()

  byteValue, err := ioutil.ReadAll(jsonFile)
  if err != nil {
    return nil, err
  }

  var events []Event
  err = json.Unmarshal(byteValue, &events)
  return events, err
}

func processActions(actions []Action) {
  for _, action := range actions {
    if len(action.Parameters) > 0 {
      processParameter(action.Parameters[0])
      if len(action.Parameters) > 1 {
        if translatedParam, err := translateParameter(action.Parameters[1], "th"); err != nil {
          fmt.Println("Translation error:", err)
        } else {
          action.Parameters[1] = translatedParam
          fmt.Println("Translated Parameter:", action.Parameters[1])
        }
      }
    } else {
      fmt.Println("Action has no parameters")
    }
  }
}

func processParameter(param interface{}) {
  switch v := param.(type) {
  case string:
    if translatedParam, err := translateText(v, "th"); err != nil {
      fmt.Println("Translation error:", err)
    } else {
      fmt.Println("Translated Parameter:", translatedParam)
    }
  case int:
    fmt.Println("Integer Parameter:", v)
  default:
    fmt.Println("Other Parameter Type:", v)
  }
}

func translateText(text string, targetLanguage string) (string, error) {
  t := translator.New()
  result, err := t.Translate(text, "auto", targetLanguage)
  if err != nil {
    return "", err
  }
  return result.Text, nil
}

func translateParameter(param interface{}, targetLanguage string) (interface{}, error) {
  switch v := param.(type) {
  case string:
    if translatedText, err := translateText(v, targetLanguage); err != nil {
      return "", err
    } else {
      return translatedText, nil
    }
  case int, float64:
    return v, nil
  default:
    return "", fmt.Errorf("unsupported parameter type: %T", v)
  }
}

func exportTranslatedJSON(events []Event, originalFilePath string) error {
  // Extract the filename from the original file path
  originalFilename := filepath.Base(originalFilePath)
  // Replace the file extension with "Translated.json"
  translatedFilename := strings.TrimSuffix(originalFilename, filepath.Ext(originalFilename)) + ".json"
  translatedFilePath := filepath.Join(filepath.Dir(originalFilePath), translatedFilename)

  translatedEvents := make([]Event, len(events))
  for i, event := range events {
    translatedEvent := Event{
      ID:       event.ID,
      Name:     event.Name,
      List:     event.List,
      SwitchID: event.SwitchID,
      Trigger:  event.Trigger,
    }

    translatedName, err := translateText(event.Name, "th")
    if err != nil {
      return err
    }
    translatedEvent.Name = translatedName

    for _, action := range translatedEvent.List {
      if len(action.Parameters) > 1 {
        if translatedParam, err := translateParameter(action.Parameters[1], "th"); err == nil {
          action.Parameters[1] = translatedParam
        }
      }
    }

    translatedEvents[i] = translatedEvent
  }

  // Encode translatedEvents to JSON
  jsonBytes, err := json.Marshal(translatedEvents)
  if err != nil {
    return err
  }

  // Write JSON to file
  err = ioutil.WriteFile(translatedFilePath, jsonBytes, 0644)
  if err != nil {
    return err
  }

  return nil
}
