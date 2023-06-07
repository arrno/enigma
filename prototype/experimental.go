package enigma

import (
	"errors"
	"fmt"
	"reflect"
)

func insertPathPtr(path []string, data any, newValue any) error {

	if data == nil {
		return errors.New("Invalid path.")
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
				if len(path) == 1 && (val.Type().Elem() == reflect.TypeOf(newValue) || val.Type().Elem().Kind() == reflect.Interface) {
					val.SetMapIndex(k, reflect.ValueOf(newValue))
					return nil
				} else if len(path) == 1 {
					return errors.New("Unable to set value.")
				} else {
					// This doesn't work
					next := val.MapIndex(k).Addr().Interface()
					return insertPathPtr(path[1:], next, newValue)
				}
			}
		}
		return errors.New("Not found.")
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			if fmt.Sprintf("%d", i) == path[0] {
				if len(path) == 1 && (val.Type().Elem() == reflect.TypeOf(newValue) || val.Type().Elem().Kind() == reflect.Interface) {
					val.Index(i).Set(reflect.ValueOf(newValue))
					return nil
				} else if len(path) == 1 {
					return errors.New("Unable to set value.")
				} else {
					// This doesn't work
					next := val.Index(i).Addr().Interface()
					return insertPathPtr(path[1:], next, newValue)
				}
			}
		}
		return errors.New("Not found.")
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if !val.Type().Field(i).IsExported() {
				continue
			}
			if val.Type().Field(i).Name == path[0] {
				if len(path) == 1 && (val.Field(i).CanSet() &&
					val.Field(i).Type() == reflect.TypeOf(newValue) || val.Field(i).Type().Kind() == reflect.Interface) {
					val.Field(i).Set(reflect.ValueOf(newValue))
					return nil
				} else if len(path) == 1 {
					return errors.New("Unable to set value.")
				} else {
					next := val.Field(i).Addr().Interface()
					return insertPathPtr(path[1:], next, newValue)
				}
			}
		}
		return errors.New("Not found.")
	default:
		return errors.New("Invalid node type in path.")
	}

}
