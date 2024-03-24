package main

import (
    "fmt"
    "io/ioutil"
    "os"

    "github.com/minio/simdjson-go"
    "github.com/Conight/go-googletrans"
)


func saveJSON(jsonData interface{}, filename string) error {
  data, err := json.MarshalIndent(jsonData, "", "  ")
  if err != nil {
    return err
  }
  err = ioutil.WriteFile(filename, data, 0644)
  if err != nil {
    return err
  }
  return nil
}

func readJSON(filename string) (interface{}, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    p := simdjson.NewParser()
    parsed, err := p.Parse(data, nil)
    if err != nil {
        return nil, err
    }

    var jsonData interface{}
    err = parsed.Unmarshal(&jsonData)
    if err != nil {
        return nil, err
    }

    return jsonData, nil
}
func translateAll(jsonData interface{}, targetLang string) interface{} {
  switch jsonData.(type) {
  case string:
    translatedText, err := translateText(jsonData.(string), targetLang)
    if err != nil {
      fmt.Println("Error translating text:", err)
      os.Exit(1)
    }
    return translatedText
  case map[string]interface{}:
    translatedMap := make(map[string]interface{})
    for key, value := range jsonData.(map[string]interface{}) {
      translatedMap[key] = translateAll(value, targetLang)
    }
    return translatedMap
  case []interface{}:
    translatedSlice := make([]interface{}, len(jsonData.([]interface{})))
    for i, value := range jsonData.([]interface{}) {
      translatedSlice[i] = translateAll(value, targetLang)
    }
    return translatedSlice
  default:
    return jsonData
  }
}



func translateText(text string, targetLang string) (string, error) {
  t := translator.New()
  result, err := t.Translate(text, targetLang, "auto")
  if err != nil {
    return "", err
  }
  return result.Text, nil
}




func main() {
   filename := "your_json_file.json"
targetLang := "en"

jsonData, err := readJSON(filename)
if err != nil {
  fmt.Println("Error reading JSON:", err)
  os.Exit(1)
}

translatedData := translateAll(jsonData, targetLang)

err = saveJSON(translatedData, "translated.json")
if err != nil {
  fmt.Println("Error saving translated JSON:", err)
  os.Exit(1)
}

fmt.Println("Translation completed and saved to translated.json")

}
