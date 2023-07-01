package enigma

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// walk will walk through any unknown structure of unknown depth/type.
func walk(data any, depth uint) {
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
			walk(val.MapIndex(k).Interface(), depth+1)
		}
		break
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			walk(val.Index(i).Interface(), depth+1)
		}
		break
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if !val.Type().Field(i).IsExported() {
				continue
			}
			walk(val.Field(i).Interface(), depth+1)
		}
		break
	default:
	}
}

// display will pretty print the data.
func display(data any) {
	r, _ := json.MarshalIndent(data, "", "    ")
	fmt.Println(string(r))
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

	wasPointer := false
	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Pointer {
		wasPointer = true
		val = val.Elem()
	}

	handlePtr := func() any {
		if wasPointer && val.CanAddr() {
			return val.Addr().Interface()
		}
		return val.Interface()
	}

	switch t := val.Type().Kind(); t {
	case reflect.Map:
		for _, k := range val.MapKeys() {
			s, ok := k.Interface().(string)
			if (ok && s == path[0]) || (fmt.Sprintf("%d", k.Interface()) == path[0]) {
				next := val.MapIndex(k).Interface()
				if r, err := insertPath(path[1:], next, newValue); err != nil {
					return handlePtr(), err
				} else if val.Type().Elem() == reflect.TypeOf(r) || val.Type().Elem().Kind() == reflect.Interface {
					val.SetMapIndex(k, reflect.ValueOf(r))
					return handlePtr(), nil
				} else {
					return handlePtr(), errors.New("Mismatched type.")
				}
			}
		}
		break
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			if fmt.Sprintf("%d", i) == path[0] {
				next := val.Index(i).Interface()
				if r, err := insertPath(path[1:], next, newValue); err != nil {
					return handlePtr(), err
				} else if val.Type().Elem() == reflect.TypeOf(r) || val.Type().Elem().Kind() == reflect.Interface {
					val.Index(i).Set(reflect.ValueOf(r))
					return handlePtr(), nil
				} else {
					return handlePtr(), errors.New("Mismatched type.")
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
					return handlePtr(), err
				} else if val.Field(i).CanSet() &&
					val.Field(i).Type() == reflect.TypeOf(r) || val.Field(i).Type().Kind() == reflect.Interface {
					val.Field(i).Set(reflect.ValueOf(r))
					return handlePtr(), nil
				} else {
					return handlePtr(), errors.New("Mismatched type.")
				}
			}
		}
		break
	default:
		return nil, errors.New("Invalid node type in path.")
	}

	return nil, errors.New("Not found.")
}

// dropPath prunes the data at the target path.
// 
// UNTESTED
func dropPath(path []string, data any) (any, error) {

	if data == nil {
		return data, errors.New("Invalid path.")
	} else if len(path) == 0 {
		return nil, errors.New("Not found.")
	}

	wasPointer := false
	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Pointer {
		wasPointer = true
		val = val.Elem()
	}

	handlePtr := func() any {
		if wasPointer && val.CanAddr() {
			return val.Addr().Interface()
		}
		return val.Interface()
	}

	switch t := val.Type().Kind(); t {
	case reflect.Map:
		for _, k := range val.MapKeys() {
			s, ok := k.Interface().(string)
			if !ok {
				s = fmt.Sprintf("%d", k.Interface())
			}
			if s == path[0] {
				if len(path) != 1 {
					next := val.MapIndex(k).Interface()
					if r, err := dropPath(path[1:], next); err != nil {
						return handlePtr(), err
					} else if val.Type().Elem() == reflect.TypeOf(r) || val.Type().Elem().Kind() == reflect.Interface {
						val.SetMapIndex(k, reflect.ValueOf(r))
						return handlePtr(), nil
					} else {
						return handlePtr(), errors.New("Mismatched type.")
					}
				} else {
					val = dropMapKey(val, s)
					return handlePtr(), nil
				}
			}
		}
		break
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			if fmt.Sprintf("%d", i) == path[0] {
				if len(path) != 1 {
					next := val.Index(i).Interface()
					if r, err := dropPath(path[1:], next); err != nil {
						return handlePtr(), err
					} else if val.Type().Elem() == reflect.TypeOf(r) || val.Type().Elem().Kind() == reflect.Interface {
						val.Index(i).Set(reflect.ValueOf(r))
						return handlePtr(), nil
					} else {
						return handlePtr(), errors.New("Mismatched type.")
					}
				} else {
					val = dropSliceIndex(val, i)
					return handlePtr(), nil
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
				if len(path) != 1 {
					next := val.Field(i).Interface()
					if r, err := dropPath(path[1:], next); err != nil {
						return handlePtr(), err
					} else if val.Field(i).CanSet() &&
						val.Field(i).Type() == reflect.TypeOf(r) || val.Field(i).Type().Kind() == reflect.Interface {
						val.Field(i).Set(reflect.ValueOf(r))
						return handlePtr(), nil
					} else {
						return handlePtr(), errors.New("Mismatched type.")
					}
				} else if wasPointer {
					val = zeroSliceField(reflect.ValueOf(data), val.Type().Field(i).Name)
				}
				return handlePtr(), nil
			}
		}
		break
	default:
		return nil, errors.New("Invalid node type in path.")
	}

	return nil, errors.New("Not found.")
}

func dropSliceIndex(val reflect.Value, index int) reflect.Value {
	newVal := reflect.MakeSlice(val.Type(), 0, val.Cap() -1)
	for i := 0; i < val.Len(); i++ {
		if i == index {
			continue
		}
		newVal = reflect.Append(newVal, val.Index(i))
	}
	return newVal
}

func dropMapKey(val reflect.Value, key string) reflect.Value {
	newVal := reflect.MakeMapWithSize(val.Type(), val.Len() - 1)
	for _, k := range val.MapKeys() {
		s, ok := k.Interface().(string)
		if (ok && s == key) || (fmt.Sprintf("%d", k.Interface()) == key) {
			continue
		}
		newVal.SetMapIndex(k, val.MapIndex(k))
	}
	return newVal
}

func zeroSliceField(val reflect.Value, field string) reflect.Value {
	wasPtr := false
	if val.Kind() == reflect.Pointer {
		wasPtr = true
		val = val.Elem()
	}
	for i := 0; i < val.NumField(); i++ {
		if !val.Type().Field(i).IsExported() {
			continue
		}
		if val.Type().Field(i).Name == field && val.Field(i).CanSet(){
			z := reflect.Zero(val.Field(i).Type())
			val.Field(i).Set(z)
		}
	}
	if wasPtr {
		return val.Addr()
	}
	return val
}