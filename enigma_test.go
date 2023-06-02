package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertPath(t *testing.T) {
	// TODO handle panic on set type mismatch
	paths := [][]string{
		{"hi", "1"},
		{"foo"},
		{"bar", "buz", "4"},
		{"biz", "box", "fix", "1"},
		{"biz", "box", "mix"},
		{"not", "found"},                // should not work
		{"fac", "slic", "0", "private"}, // should not work
		// wrong types TODO
		// {"strict","0"},
		// {"fac", "slic", "0", "WidgetSize"},
		{"fac", "slic", "0", "Gadgets", "0", "Name"},
		{"fac", "slic", "0", "WidgetName"}, // should not work
		{"fac", "nop", "WidgetColor"},      // should not work
		{"ptr", "0", "0", "0"},
	}
	for _, p := range paths {
		insertPath(p, &data, "REPLACED")
	}
	result, _ := json.MarshalIndent(data, "", "    ")
	assert.Equal(t, replaced, result)

}
func TestQueryValue(t *testing.T) {
	actual := []string{}
	paths := []string{}
	queryValue([]string{}, data, nil, &paths)
	actual = append(actual, paths...)
	paths = []string{}
	queryValue([]string{}, data, true, &paths)
	actual = append(actual, paths...)
	paths = []string{}
	queryValue([]string{}, data, "shhh", &paths)
	actual = append(actual, paths...)
	paths = []string{}
	queryValue([]string{}, data, "blue", &paths)
	actual = append(actual, paths...)
	paths = []string{}
	queryValue([]string{}, data, "p", &paths)
	actual = append(actual, paths...)
	expected := []string{
		"foo",
		"bar.buz.3",
		"bar.buz.5",
		"biz.box.fox",
		"biz.box.fix.3",
		"biz.box.fix.5.fiz",
		"biz.double.sep.double",
		"fac.slic.0.Gadgets.0.Fidgety",
		"fac.slic.0.WidgetColor",
		"ptr.0.0.0",
	}
	// Sort?
	assert.Equal(t, len(expected), len(actual))
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], actual[i])
	}
}

func TestQueryKey(t *testing.T) {
	actual := []string{}
	paths := []string{}
	queryKey([]string{}, data, "double", &paths)
	actual = append(actual, paths...)
	paths = []string{}
	queryKey([]string{}, data, "box", &paths)
	actual = append(actual, paths...)
	paths = []string{}
	queryKey([]string{}, data, "0", &paths)
	actual = append(actual, paths...)
	expected := []string{
		"biz.box.double",
		"biz.double",
		"biz.double.sep.double",
		"biz.box",
		"ptr.0",
		"ptr.0.0",
		"ptr.0.0.0",
		"fac.wid.Gadgets.0",
		"fac.slic.0",
		"fac.slic.0.Gadgets.0",
		"strict.0",
		"bar.buz.0",
		"biz.box.fix.0",
	}
	// Sort?
	assert.Equal(t, len(expected), len(actual))
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], actual[i])
	}
}
func TestQueryPath(t *testing.T) {
	paths := [][]string{
		{"hi"},
		{"hi", "1"},
		{"foo"},
		{"bar", "buz", "4"},
		{"biz", "box", "fix", "1"},
		{"biz", "box", "mix"},
		{"not", "found"},
		{"fac", "slic", "0", "private"},
		{"fac", "slic", "0", "WidgetColor"},
		{"fac", "slic", "0", "Gadgets", "0", "Name"},
		{"ptr", "0", "0", "0"},
	}
	actual := []any{}
	for _, p := range paths {
		val, err := queryPath(p, data)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			actual = append(actual, val)
		}
	}
	expected := []any{
		map[int]any{1: "hello"},
		"hello",
		nil,
		7,
		"fun",
		"pop",
		"blue",
		"doop",
		"p",
	}
	assert.Equal(t, len(expected), len(actual))
	// reflect deep equal?
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], actual[i])
	}
}

// <--------------------------------------- Fake data ------------------------------------------>

type Widget struct {
	WidgetColor string
	WidgetSize  int
	Gadgets     []any
	private     string
}

type Fidget struct {
	Fidgety bool
	Name    string
}

var data = map[string]any{
	"strict": []bool{false},
	// generic
	"hi":  map[int]any{1: "hello"},
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
				"fiz": nil,
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
	"ptr": &[]any{&[]any{&[]any{"p"}}},
	// with structs
	"fac": map[string]any{
		"wid": &Widget{
			private:     "shhh",
			WidgetColor: "red",
			WidgetSize:  10,
			Gadgets:     []any{"a", "b", "c"},
		},
		"nop": Widget{
			private:     "no",
			WidgetColor: "yellow",
			WidgetSize:  12,
		},
		"slic": []any{
			&Widget{
				private:     "quiet",
				WidgetColor: "blue",
				WidgetSize:  20,
				Gadgets: []any{
					&Fidget{
						Fidgety: true,
						Name:    "doop",
					},
				},
			},
		},
	},
}

var replaced string = `{
    "bar": {
        "buz": [
            0,
            "1",
            2,
            null,
            "REPLACED",
            null
        ]
    },
    "biz": {
        "box": {
            "double": false,
            "fix": [
                7,
                "REPLACED",
                9,
                null,
                10,
                {
                    "fiz": null
                }
            ],
            "fox": null,
            "mix": "REPLACED"
        },
        "double": {
            "sep": {
                "double": true
            }
        }
    },
    "fac": {
        "nop": {
            "WidgetColor": "yellow",
            "WidgetSize": 12,
            "Gadgets": null
        },
        "slic": [
            {
                "WidgetColor": "blue",
                "WidgetSize": 20,
                "Gadgets": [
                    {
                        "Fidgety": true,
                        "Name": "REPLACED"
                    }
                ]
            }
        ],
        "wid": {
            "WidgetColor": "red",
            "WidgetSize": 10,
            "Gadgets": [
                "a",
                "b",
                "c"
            ]
        }
    },
    "foo": "REPLACED",
    "hi": {
        "1": "REPLACED"
    },
    "ptr": [
        [
            [
                "REPLACED"
            ]
        ]
    ],
    "strict": [
        false
    ]
}`
