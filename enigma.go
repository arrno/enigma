package enigma

import (
	"errors"
	"reflect"
	"strings"
)

// <--------------------------------------- Pub API ------------------------------------------->
// <------------------------------------------------------------------------------------------->

var SUPPORTEDTYPES = []reflect.Kind{reflect.Map, reflect.Slice, reflect.Struct}

type Enigma struct {
	data any
}

// NewEnigma takes a pointer to a struct, map, or slice and returns a pointer to a new Enigma.
func NewEnigma(data any) (*Enigma, error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Pointer {
		return nil, errors.New("Must pass pointer to struct, slice, or map.")
	}
	for i, k := range SUPPORTEDTYPES {
		if v.Elem().Kind() == k {
			break
		}
		if i == len(SUPPORTEDTYPES)-1 {
			return nil, errors.New("Must pass pointer to struct, slice, or map.")
		}
	}
	e := new(Enigma)
	e.data = data
	return e, nil
}

// Display pretty prints the underlying data.
func (e *Enigma) Display() {
	display(e.data)
}

// QueryValue returns all paths that lead to an instance of the target value.
func (e *Enigma) QueryValue(value any) []string {
	results := []string{}
	queryValue([]string{}, e.data, value, &results)
	return results
}

// QueryKey returns all paths that lead to an instance of the target key.
//
// A key is either a map key, a slice index, or a struct field in string format.
func (e *Enigma) QueryKey(key string) []string {
	results := []string{}
	queryKey([]string{}, e.data, key, &results)
	return results
}

// QueryPath returns the value located at the provided path if the path exists.
func (e *Enigma) QueryPath(path string) (any, error) {
	rawPath := strings.Split(path, ".")
	return queryPath(rawPath, e.data)
}

// InsertByPath inserts the value at the provided path if the path exists.
func (e *Enigma) InsertByPath(path string, value any) error {
	rawPath := strings.Split(path, ".")
	_, err := insertPath(rawPath, e.data, value)
	return err
}

// InsertByValue replaces all instances of 'find' with 'replace'.
func (e *Enigma) InsertByValue(find any, replace any) (err error) {
	paths := e.QueryValue(find)
	for _, path := range paths {
		err = e.InsertByPath(path, replace)
	}
	return err
}

// InsertByKey updates any instances of the target key to hold the new value.
func (e *Enigma) InsertByKey(key string, replace any) (err error) {
	paths := e.QueryKey(key)
	for _, path := range paths {
		err = e.InsertByPath(path, replace)
	}
	return err
}

// DropByPath prunes the data at the target path.
func (e *Enigma) DropByPath(path string) error {
	rawPath := strings.Split(path, ".")
	_, err := dropPath(rawPath, e.data)
	return err
}

// DropByValue prunes all parent nodes that hold a target value.
func (e *Enigma) DropByValue(value any) (err error) {
	paths := e.QueryValue(value)
	for _, path := range paths {
		err = e.DropByPath(path)
	}
	return err
}

// Data is a getter for Enigma data.
func (e *Enigma) Data() any {
	return e.data
}
