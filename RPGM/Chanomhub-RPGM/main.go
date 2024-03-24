	package main

	import (
		"fmt"
		gt "github.com/bas24/googletranslatefree"
	)

	func main(){
		const text string = `Hello, World!`
		// you can use "auto" for source language
		// so, translator will detect language
		result, _ := gt.Translate(text, "en", "es")
		fmt.Println(result)
		// Output: "Hola, Mundo!"
	}
