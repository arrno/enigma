package enigma

import (
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Str struct{
	One string
	Two string
	Three int
	Four int
}
type StrStr struct{
	Control string
	Str *Str
}

func TestRun(t *testing.T) {
	// run in order
	QueryValue(t)
	QueryKey(t)
	QueryPath(t)
	InsertPath(t)
}

func TestDrop(t *testing.T) {

	a := []string{"one","two","three","four"}
	a = dropSliceIndex(reflect.ValueOf(a), 1).Interface().([]string)
	assert.Equal(t, a, []string{"one","three","four"})
	a = []string{"one","two","three","four"}
	a = dropSliceIndex(reflect.ValueOf(a), 3).Interface().([]string)
	assert.Equal(t, a, []string{"one","two","three"})

	m := map[string]any{
		"one":true,
		"two":false,
		"three":true,
		"four":false,
	}
	m = dropMapKey(reflect.ValueOf(m), "two").Interface().(map[string]any)
	assert.True(t, reflect.DeepEqual(m, map[string]any{"one":true,"three":true,"four":false}))
	m["two"] = false
	m = dropMapKey(reflect.ValueOf(m), "three").Interface().(map[string]any)
	assert.True(t, reflect.DeepEqual(m, map[string]any{"one":true,"two":false,"four":false}))

	s := &Str{"One", "Two", 5, 7}
	s = zeroSliceField(reflect.ValueOf(s), "One").Interface().(*Str)
	assert.Equal(t, []any{s.One, s.Two, s.Three, s.Four}, []any{"", "Two", 5, 7})
	s.One = "One"
	s = zeroSliceField(reflect.ValueOf(s), "Four").Interface().(*Str)
	assert.Equal(t, []any{s.One, s.Two, s.Three, s.Four}, []any{"One", "Two", 5, 0})
}

func makeData() map[string]any {
	data := map[string]any{
		"one":[]any{0,1,2,[]int{}},
		"two":map[string]any{"a":0,"b": &Str{"One", "Two", 5, 7}},
		"three": &StrStr{"Blue", &Str{"ZOne", "ZTwo", 50, 70}},
	}
	return data
}
func TestDropRecursive(t *testing.T) {

	data := makeData()
	dropPath([]string{"one","3"}, &data)
	assert.True(t, reflect.DeepEqual(
		data,
		map[string]any{
			"one":[]any{0,1,2},
			"two":map[string]any{"a":0,"b": &Str{"One", "Two", 5, 7}},
			"three": &StrStr{"Blue", &Str{"ZOne", "ZTwo", 50, 70}},
		},
	))

	data = makeData()
	dropPath([]string{"two"}, &data)
	assert.True(t, reflect.DeepEqual(
		data,
		map[string]any{
			"one":[]any{0,1,2,[]int{}},
			"three": &StrStr{"Blue", &Str{"ZOne", "ZTwo", 50, 70}},
		},
	))

	data = makeData()
	dropPath([]string{"two"}, &data)
	assert.True(t, reflect.DeepEqual(
		data,
		map[string]any{
			"one":[]any{0,1,2,[]int{}},
			"three": &StrStr{"Blue", &Str{"ZOne", "ZTwo", 50, 70}},
		},
	))

	data = makeData()
	dropPath([]string{"two", "a"}, &data)
	assert.True(t, reflect.DeepEqual(
		data,
		map[string]any{
			"one":[]any{0,1,2,[]int{}},
			"two":map[string]any{"b": &Str{"One", "Two", 5, 7}},
			"three": &StrStr{"Blue", &Str{"ZOne", "ZTwo", 50, 70}},
		},
	))

	data = makeData()
	dropPath([]string{"two", "b"}, &data)
	assert.True(t, reflect.DeepEqual(
		data,
		map[string]any{
			"one":[]any{0,1,2,[]int{}},
			"two":map[string]any{"a":0},
			"three": &StrStr{"Blue", &Str{"ZOne", "ZTwo", 50, 70}},
		},
	))

	data = makeData()
	dropPath([]string{"two", "b", "Three"}, &data)
	assert.True(t, reflect.DeepEqual(
		data,
		map[string]any{
			"one":[]any{0,1,2,[]int{}},
			"two":map[string]any{"a":0,"b": &Str{"One", "Two", 0, 7}},
			"three": &StrStr{"Blue", &Str{"ZOne", "ZTwo", 50, 70}},
		},
	))

	data = makeData()
	dropPath([]string{"three","Str"}, &data)
	assert.True(t, reflect.DeepEqual(
		data,
		map[string]any{
			"one":[]any{0,1,2,[]int{}},
			"two":map[string]any{"a":0,"b": &Str{"One", "Two", 5, 7}},
			"three": &StrStr{Control: "Blue", Str: nil},
		},
	))

	data = makeData()
	dropPath([]string{"three", "Str", "One"}, &data)
	assert.True(t, reflect.DeepEqual(
		data,
		map[string]any{
			"one":[]any{0,1,2,[]int{}},
			"two":map[string]any{"a":0,"b": &Str{"One", "Two", 5, 7}},
			"three": &StrStr{"Blue", &Str{"", "ZTwo", 50, 70}},
		},
	))

	// root list pointer...
	ldata := []any{0,1,2}
	dropPath([]string{"0"}, &ldata)
	assert.Equal(t, ldata, []any{1,2})

	ldata = []any{0,1,map[string]int{"a":1,"b":2}}
	dropPath([]string{"2","b"}, &ldata)
	assert.True(t, reflect.DeepEqual(ldata, []any{0,1,map[string]int{"a":1}}))

	// root struct pointer...
	sdata := StrStr{
		Control: "Blue",
		Str: &Str{"One","Two",5,7},
	}
	dropPath([]string{"Str"}, &sdata)
	assert.True(t, reflect.DeepEqual(sdata, StrStr{Control: "Blue",Str: nil}))
	sdata.Str = &Str{"One","Two",5,7}
	dropPath([]string{"Str","One"}, &sdata)
	assert.True(t, reflect.DeepEqual(sdata, StrStr{Control: "Blue",Str: &Str{"","Two",5,7}}))

}

func InsertPath(t *testing.T) {
	type testCase struct {
		set  any
		path []string
		fail bool
	}
	tests := []testCase{
		{
			set:  "REPLACE",
			path: []string{"hi", "1"},
		},
		{
			set:  "REPLACE",
			path: []string{"foo"},
		},
		{
			set:  "REDO",
			path: []string{"foo"},
		},
		{
			set:  "REPLACE",
			path: []string{"bar", "buz", "4"},
		},
		{
			set:  "REPLACE",
			path: []string{"biz", "box", "fix", "1"},
		},
		{
			set:  "REPLACE",
			path: []string{"biz", "box", "mix"},
		},
		{
			set:  "REPLACE",
			path: []string{"not", "found"},
			fail: true,
		},
		{
			set:  "REPLACE",
			path: []string{"fac", "slic", "0", "private"},
			fail: true,
		},
		{
			set:  "REPLACE",
			path: []string{"strict", "0"},
			fail: true,
		},
		{
			set:  "REPLACE",
			path: []string{"strict", "0", "WidgetSize"},
			fail: true,
		},
		{
			set:  "REPLACE",
			path: []string{"fac", "slic", "0", "Gadgets", "0", "Name"},
		},
		{
			set:  "REDO",
			path: []string{"fac", "slic", "0", "Gadgets", "0", "Name"},
		},
		{
			set:  "REPLACE",
			path: []string{"fac", "slic", "0", "WidgetName"},
			fail: true,
		},
		{
			set:  "REPLACE",
			path: []string{"fac", "nop", "WidgetColor"},
			fail: true,
		},
		{
			set:  "REPLACE",
			path: []string{"ptr", "0", "0", "0"},
		},
		{
			set:  "REDO",
			path: []string{"ptr", "0", "0", "0"},
		},
	}
	for _, tst := range tests {
		_, err := insertPath(tst.path, &data, tst.set)
		if !tst.fail {
			assert.Nil(t, err)
			actual, _ := queryPath(tst.path, data)
			assert.Equal(t, tst.set, actual)
		} else {
			assert.NotNil(t, err)
		}
	}
}

func TestInsert(t *testing.T) {
	// insert typed map into typed slice
	d := []map[string]int{{}}
	insertPath([]string{"0"}, d, map[string]int{"foo": 7})
	actual, _ := queryPath([]string{"0", "foo"}, d)
	assert.Equal(t, 7, actual)

	// insert into typed map
	insertPath([]string{"0", "foo"}, d, 9)
	actual, _ = queryPath([]string{"0", "foo"}, d)
	assert.Equal(t, 9, actual)

	// insert a struct into a typed map
	sd := map[string]Fidget{"one": {}}
	insertPath([]string{"one"}, sd, Fidget{Name: "Sassy"})
	actual, _ = queryPath([]string{"one", "Name"}, sd)
	assert.Equal(t, "Sassy", actual)

	// insert a slice into a typed map
	sld := map[string][]int{"one": {}}
	insertPath([]string{"one"}, sld, []int{0, 1, 2})
	actual, _ = queryPath([]string{"one", "2"}, sld)
	assert.Equal(t, 2, actual)

	// write to something beyond a non pointer struct
	sw := map[string]Widget{"one": {Gadgets: []any{"foo"}}}
	insertPath([]string{"one", "Gadgets", "0"}, sw, "bar")
	actual, _ = queryPath([]string{"one", "Gadgets", "0"}, sw)
	assert.Equal(t, "bar", actual)

	// should fail because writing to non ptr struct
	sw = map[string]Widget{"one": {Gadgets: []any{"foo"}}}
	insertPath([]string{"one", "Gadgets"}, sw, []any{"bar"})
	actual, _ = queryPath([]string{"one", "Gadgets", "0"}, sw)
	assert.Equal(t, "foo", actual)

	swp := map[string]*Widget{"one": {Gadgets: []any{"foo"}}}
	insertPath([]string{"one", "Gadgets", "0"}, swp, "bar")
	actual, _ = queryPath([]string{"one", "Gadgets", "0"}, swp)
	assert.Equal(t, "bar", actual)

}

func QueryValue(t *testing.T) {
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
	sort.Strings(expected)
	sort.Strings(actual)
	assert.Equal(t, len(expected), len(actual))
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], actual[i])
	}
}

func QueryKey(t *testing.T) {
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
	sort.Strings(expected)
	sort.Strings(actual)
	// Sort?
	assert.Equal(t, len(expected), len(actual))
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], actual[i])
	}
}

func QueryPath(t *testing.T) {
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
			// fmt.Println(err.Error())
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
