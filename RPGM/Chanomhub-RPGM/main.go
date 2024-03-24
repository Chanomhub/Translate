package main

import (
	"encoding/json"
	"fmt"
	"flag"
	gt "github.com/bas24/googletranslatefree"
	"io/ioutil"
	"os"
)

func main() {
	// ตั้งค่า flag
	inputPath := flag.String("input", "", "Path to the input JSON file")
	outputPath := flag.String("output", "", "Path to the output JSON file")
	sourceLang := flag.String("source", "auto", "Source language code")
	targetLang := flag.String("target", "en", "Target language code")
	flag.Parse()

	// ตรวจสอบ flag
	if *inputPath == "" || *outputPath == "" {
		fmt.Println("Usage: go run translate.go -input <input.json> -output <output.json> -source <source_lang> -target <target_lang>")
		return
	}

	// อ่านไฟล์ JSON
	data, err := ioutil.ReadFile(*inputPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// แปลง JSON เป็น map
	var jsonData map[string]interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		fmt.Println(err)
		return
	}

	// ดึงข้อมูล name
	name, ok := jsonData["name"]
	if !ok {
		fmt.Println("No 'name' field found in the JSON file")
		return
	}

	// แปลภาษา
	translatedText, err := gt.Translate(name.(string), *sourceLang, *targetLang)
	if err != nil {
		fmt.Println(err)
		return
	}

	// แทรกข้อความแปลกลับไปยัง .json
	jsonData["translated_name"] = translatedText

	// เขียน JSON ไปยังไฟล์
	output, err := json.MarshalIndent(jsonData, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(*outputPath, output, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Translation successful!")
}
