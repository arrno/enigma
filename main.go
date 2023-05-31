package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	test3()
}

func test3(){
	// data := map[string]any{
	// 	"foo": map[string]any{
	// 		"bar": "biz",
	// 	},
	// }
	paths := [][]string{
		// {"hi", "1"},
		{"foo"},
		// {"bar","buz","4"},
		// {"biz", "box", "fux", "1"},
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

func test1(){
	paths := []string{}
	queryValue([]string{}, data, nil, &paths)
	for _, p := range paths {
		fmt.Println(p)
	}
}

func test2() {
	paths := [][]string{
		{"hi"},
		{"hi", "1"},
		{"foo"},
		{"bar","buz","4"},
		{"biz", "box", "fux", "1"},
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

var data = map[string]any{
	"hi": map[int]any{ 1: "hello"},
	"foo": nil,
	"bar": map[string]any{
		"buz": []any{0, "1", 2, nil, 7, nil},
	},
	"biz": map[string]any{
		"box": map[string]any{
			"mix": "pop",
			"fox": nil,
			"fux": []any{7, "fun", 9, nil, 10, map[string]any{
					"fiz":nil,
				},
			},
		},
	},
}