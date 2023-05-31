package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func walk(data any) {
	val := reflect.ValueOf(data)
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
			walk(val.Field(i).Interface())
		}
	}
}

// get paths by value
func queryValue(seen []string, data any, target any, results *[]string) {
	if reflect.DeepEqual(data, target) {
		*results = append(*results, strings.Join(seen, "."))
		return
	}
	val := reflect.ValueOf(data)
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
			seen = append(seen, val.Type().Field(i).Name)
			queryValue(seen, val.Field(i).Interface(), target, results)
		}
		break
	}
}

// get value by path
func queryPath(path []string, data any) (any, error) {
	if len(path) == 0 {
		return data, nil
	}
	val := reflect.ValueOf(data)
	switch t := val.Type().Kind(); t {
	case reflect.Map:
        for _, k := range val.MapKeys() {
			s, ok := k.Interface().(string)
			if (ok && s == path[0]) || (fmt.Sprintf("%d", k.Interface()) == path[0]){
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

// insert value at path
func insertPath(path []string, data any, newValue any) (any, error) {
	
	if len(path) == 0 {
		return newValue, nil
	}

	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	switch t := val.Type().Kind(); t {
	case reflect.Map:
        for _, k := range val.MapKeys() {
			s, ok := k.Interface().(string)
			if (ok && s == path[0]) || (fmt.Sprintf("%d", k.Interface()) == path[0]){
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
			if val.Type().Field(i).Name == path[0] {
				next := val.Field(i).Interface()
				if r, err := insertPath(path[1:], &next, newValue); err != nil {
					return val.Interface(), err
				} else {
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