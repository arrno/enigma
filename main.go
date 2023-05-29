package main

import "fmt"

func main() {
	data := map[string]any{
		"foo": "bar",
	}
	enigma, _ := NewEnigma(data)
	fmt.Println(enigma.Data())
}
