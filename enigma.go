package main

import (
	"errors"
	"reflect"
	"strings"
)

// <--------------------------------------- Pub API ------------------------------------------->
// <------------------------------------------------------------------------------------------->

// Structs are not yet supported. Only supporting map[string]<T>
var SUPPORTEDTYPES = []reflect.Kind{reflect.Map, reflect.Slice, reflect.Struct}

type Enigma struct {
	data     any
	mapAny   *map[string]any
	sliceAny *[]any
	rootType reflect.Kind
}

// NewEnigma takes a pointer to a struct, map, or slice and returns a pointer to a new Enigma.
func NewEnigma(data any) (*Enigma, error) {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Pointer {
		return nil, errors.New("data must be a pointer to struct, map, or slice.")
	}
	e := new(Enigma)
	k := t.Elem().Kind()
	for i, v := range SUPPORTEDTYPES {
		if k == v && v == reflect.Map {
			e.rootType = v
			d := toMapAny(data)
			e.mapAny = &d
			e.data = d
			break
		} else if k == v && v == reflect.Slice {
			e.rootType = v
			d := toSliceAny(data)
			e.sliceAny = &d
			e.data = d
			break
		} else if i == len(SUPPORTEDTYPES)-1 {
			return nil, errors.New("data must be a pointer to struct, map, or slice.")
		}
	}
	return e, nil
}

// QueryValue returns all paths that lead to an instance of the target type
func (e *Enigma) QueryValue(value any) (paths []string, err error) {
	if e.data == nil {
		return paths, err
	}
	rawPaths := [][]string{}
	switch e.rootType {
	case reflect.Map:
		getValueMapPaths(*e.mapAny, []string{}, value, &rawPaths)
		break
	case reflect.Slice:
		getValueSlicePaths(*e.sliceAny, []string{}, value, &rawPaths)
		break
	case reflect.Struct:
		return paths, errors.New("Struct paths not yet supported.")
	}
	for _, path := range rawPaths {
		paths = append(paths, strings.Join(path, "."))
	}
	return paths, err
}

// QueryPath returns the value located at the provided path. If the path location
// does not exist, an error is returned.
func (e *Enigma) QueryPath(path string) (value any, err error) {
	rawPath := strings.Split(path, ".")
	switch e.rootType {
	case reflect.Map:
		value, err = getMapPathValue(*e.mapAny, rawPath)
		break
	case reflect.Slice:
		value, err = getSlicePathValue(*e.sliceAny, rawPath)
		break
	case reflect.Struct:
		return nil, errors.New("Struct paths not yet supported.")
	}
	return value, err
}

// InsertByPath inserts the value at the path location. If the location does not exist,
// an error is returned.
func (e *Enigma) InsertByPath(path string, value any) error {
	rawPath := strings.Split(path, ".")
	switch e.rootType {
	case reflect.Map:
		return replaceMapValues(e.mapAny, rawPath, value)
	case reflect.Slice:
		return replaceSliceValues(e.sliceAny, rawPath, value)
	case reflect.Struct:
		return errors.New("Struct paths not yet supported.")
	}
	return errors.New("Unsupported root type.")
}

// InsertByValue replaces all instances of 'find' with 'replace'.
func (e *Enigma) InsertByValue(find any, replace any) (err error) {
	paths, err := e.QueryValue(find)
	if err != nil {
		return err
	}
	for _, path := range paths {
		err = e.InsertByPath(path, replace)
	}
	return err
}

// Data is a getter for Enigma data.
func (e *Enigma) Data() any {
	return e.data
}
