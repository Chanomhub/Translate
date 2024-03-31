package main

  import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"

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

func main() {
    events, err := readJSONFile("CommonEvents.json")
    if err != nil {
        fmt.Println(err)
        return
    }

    for _, event := range events {
        fmt.Println("Original Event Name:", event.Name)

        if translatedName, err := translateText(event.Name, "th"); err != nil {
            fmt.Println("Translation error:", err)
        } else {
            event.Name = translatedName
            fmt.Println("Translated Event Name:", event.Name)
        }

        processActions(event.List)
    }

    // Export translated JSON
    if err := exportTranslatedJSON(events); err != nil {
        fmt.Println("Error exporting translated JSON:", err)
        return
    }
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


func exportTranslatedJSON(events []Event) error {
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
    err = ioutil.WriteFile("TranslatedEvents.json", jsonBytes, 0644)
    if err != nil {
        return err
    }

    return nil
}
