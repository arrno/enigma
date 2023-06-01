package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// walk will walk through any unknown structure of unknown depth/type.
func walk(data any) {
	if data == nil {
		return
	}
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	switch t := val.Type().Kind(); t {
	case reflect.Map:
		for _, k := range val.MapKeys() {
			walk(val.MapIndex(k))
		}
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			walk(val.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if !val.Type().Field(i).IsExported() {
				continue
			}
			walk(val.Field(i).Interface())
		}
	}
}

// queryKe gets paths that contain a target key.
func queryKey(seen []string, data any, target string, results *[]string) {

	if data == nil {
		return
	}
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	switch t := val.Type().Kind(); t {
	case reflect.Map:
		for _, k := range val.MapKeys() {
			scopy := make([]string, len(seen))
			copy(scopy, seen)
			s, ok := k.Interface().(string)
			if ok {
				scopy = append(scopy, s)
			} else {
				s = fmt.Sprintf("%d", k.Interface())
				scopy = append(scopy, s)
			}
			if s == target {
				*results = append(*results, strings.Join(scopy, "."))
			}
			queryKey(scopy, val.MapIndex(k).Interface(), target, results)
		}
		break
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			s := fmt.Sprintf("%d", i)
			scopy := make([]string, len(seen))
			copy(scopy, seen)
			scopy = append(scopy, s)
			if s == target {
				*results = append(*results, strings.Join(scopy, "."))
			}
			queryKey(scopy, val.Index(i).Interface(), target, results)
		}
		break
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if !val.Type().Field(i).IsExported() {
				continue
			}
			s := val.Type().Field(i).Name
			scopy := make([]string, len(seen))
			copy(scopy, seen)
			scopy = append(scopy, s)
			if s == target {
				*results = append(*results, strings.Join(scopy, "."))
			}
			queryKey(scopy, val.Field(i).Interface(), target, results)
		}
		break
	}
}

// queryValue gets all paths to a target value.
func queryValue(seen []string, data any, target any, results *[]string) {

	if reflect.DeepEqual(data, target) {
		*results = append(*results, strings.Join(seen, "."))
		return
	} else if data == nil {
		return
	}
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	switch t := val.Type().Kind(); t {
	case reflect.Map:
		for _, k := range val.MapKeys() {
			scopy := make([]string, len(seen))
			copy(scopy, seen)
			if s, ok := k.Interface().(string); ok {
				scopy = append(scopy, s)
			} else {
				scopy = append(scopy, fmt.Sprintf("%d", k.Interface()))
			}
			queryValue(scopy, val.MapIndex(k).Interface(), target, results)
		}
		break
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			scopy := make([]string, len(seen))
			copy(scopy, seen)
			scopy = append(scopy, fmt.Sprintf("%d", i))
			queryValue(scopy, val.Index(i).Interface(), target, results)
		}
		break
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if !val.Type().Field(i).IsExported() {
				continue
			}
			s := val.Type().Field(i).Name
			scopy := make([]string, len(seen))
			copy(scopy, seen)
			scopy = append(scopy, s)
			queryValue(scopy, val.Field(i).Interface(), target, results)
		}
		break
	}
}

// queryPath gets the value located at the target path.
func queryPath(path []string, data any) (any, error) {

	if len(path) == 0 {
		return data, nil
	} else if data == nil {
		return data, errors.New("Invalid path.")
	}
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	switch t := val.Type().Kind(); t {
	case reflect.Map:
		for _, k := range val.MapKeys() {
			s, ok := k.Interface().(string)
			if (ok && s == path[0]) || (fmt.Sprintf("%d", k.Interface()) == path[0]) {
				return queryPath(path[1:], val.MapIndex(k).Interface())
			}
		}
		break
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			if fmt.Sprintf("%d", i) == path[0] {
				return queryPath(path[1:], val.Index(i).Interface())
			}
		}
		break
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if !val.Type().Field(i).IsExported() {
				continue
			}
			if val.Type().Field(i).Name == path[0] {
				return queryPath(path[1:], val.Field(i).Interface())
			}
		}
		break
	default:
		return nil, errors.New("Invalid node type in path.")
	}

	return nil, errors.New("Not found.")
}

// insertPath inserts any new value at the target path location.
func insertPath(path []string, data any, newValue any) (any, error) {

	if len(path) == 0 {
		return newValue, nil
	} else if data == nil {
		return data, errors.New("Invalid path.")
	}

	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	switch t := val.Type().Kind(); t {
	case reflect.Map:
		for _, k := range val.MapKeys() {
			s, ok := k.Interface().(string)
			if (ok && s == path[0]) || (fmt.Sprintf("%d", k.Interface()) == path[0]) {
				next := val.MapIndex(k).Interface()
				if r, err := insertPath(path[1:], next, newValue); err != nil {
					return val.Interface(), err
				} else {
					val.SetMapIndex(k, reflect.ValueOf(r))
					return val.Interface(), nil
				}
			}
		}
		break
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			if fmt.Sprintf("%d", i) == path[0] {
				next := val.Index(i).Interface()
				if r, err := insertPath(path[1:], next, newValue); err != nil {
					return val.Interface(), err
				} else {
					val.Index(i).Set(reflect.ValueOf(r))
					return val.Interface(), nil
				}
			}
		}
		break
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if !val.Type().Field(i).IsExported() {
				continue
			}
			if val.Type().Field(i).Name == path[0] {
				next := val.Field(i).Interface()
				if r, err := insertPath(path[1:], next, newValue); err != nil {
					return val.Interface(), err
				} else if val.Field(i).CanSet() {
					val.Field(i).Set(reflect.ValueOf(r))
					return val.Interface(), nil
				}
			}
		}
		break
	default:
		return nil, errors.New("Invalid node type in path.")
	}

	return nil, errors.New("Not found.")
}
