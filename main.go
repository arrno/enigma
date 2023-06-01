package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	test4()
}

func test4(){
	paths := []string{}
	queryKey([]string{}, data, "double", &paths)
	for _, p := range paths {
		fmt.Println(p)
	}
	paths = []string{}
	queryKey([]string{}, data, "box", &paths)
	for _, p := range paths {
		fmt.Println(p)
	}
	paths = []string{}
	queryKey([]string{}, data, "0", &paths)
	for _, p := range paths {
		fmt.Println(p)
	}
}

// insert by path
func test3(){
	paths := [][]string{
		{"hi", "1"},
		{"foo"},
		{"bar","buz","4"},
		{"biz", "box", "fix", "1"},
		{"biz", "box", "mix"},
		{"not","found"},
	}
	// paths := [][]string{{"foo","bar"}}
	for _, p := range paths {
		_, err := insertPath(p, &data, "REPLACED")
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	js, _ := json.MarshalIndent(data, "", "    ")
	fmt.Println(string(js))
}

// query by val
func test1(){
	paths := []string{}
	queryValue([]string{}, data, nil, &paths)
	for _, p := range paths {
		fmt.Println(p)
	}
	paths = []string{}
	queryValue([]string{}, data, true, &paths)
	for _, p := range paths {
		fmt.Println(p)
	}
}

// query by path
func test2() {
	paths := [][]string{
		{"hi"},
		{"hi", "1"},
		{"foo"},
		{"bar","buz","4"},
		{"biz", "box", "fix", "1"},
		{"biz", "box", "mix"},
		{"not","found"},
	}
	for _, p := range paths {
		val, err := queryPath(p, data)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(val)
		}
	}
}

type Widget struct{
	WidgetColor string
	WidgetSize int
	Gadgets []any
	private string
}
type Fidget struct{
	Fidgety bool
	Name string
}

var data = map[string]any{
	// generic
	"hi": map[int]any{ 1: "hello"},
	"foo": nil,
	"bar": map[string]any{
		"buz": []any{0, "1", 2, nil, 7, nil},
	},
	// very nested
	"biz": map[string]any{
		"box": map[string]any{
			"mix": "pop",
			"fox": nil,
			"fix": []any{7, "fun", 9, nil, 10, map[string]any{
					"fiz":nil,
				},
			},
			"double": false,
		},
		// dup keys
		"double": map[string]any{
			"sep": map[string]any{
				"double": true,
			},
		},
	},
	// with structs
	"fac": map[string]any{
		"wid": Widget{
			private: "shhh",
			WidgetColor: "red",
			WidgetSize: 10,
			Gadgets: []any{"a", "b", "c"},
		},
		"slic": []any{
			Widget{
				private: "quiet",
				WidgetColor: "blue",
				WidgetSize: 20,
				Gadgets: []any{
					Fidget{
						Fidgety: true,
						Name: "doop",
					},
				},
			},
		},
	},
}