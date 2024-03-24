package main

import (
    "fmt"
    "io/ioutil"
    "os"

    "github.com/minio/simdjson-go"
    "github.com/Conight/go-googletrans"
)



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


func translateText(text string, targetLang string) (string, error) {
    translator := googletrans.NewTranslator()
    translatedText, err := translator.Translate(text, targetLang, "auto")
    if err != nil {
        return "", err
    }

    return translatedText, nil
}



func main() {
    filename := "your_json_file.json"
    targetLang := "en" // ตั้งค่าเป็นภาษาที่ต้องการแปลเป็น

    jsonData, err := readJSON(filename)
    if err != nil {
        fmt.Println("Error reading JSON:", err)
        os.Exit(1)
    }

    // ทำการแปลภาษาทุกข้อความใน JSON
    // (ต้องทำการ iterate ทุกๆ ข้อมูลและแปลแยกทีละข้อความ)
    translatedData := translateAll(jsonData, targetLang)

    // บันทึกผลลัพธ์ที่แปลแล้วลงในไฟล์
    // (ใช้วิธีการ serialize โครงสร้างข้อมูลกลับเป็น JSON หลังจากแปลแล้ว)
    err = saveJSON(translatedData, "translated.json")
    if err != nil {
        fmt.Println("Error saving translated JSON:", err)
        os.Exit(1)
    }

    fmt.Println("Translation completed and saved to translated.json")
}
