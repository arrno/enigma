package enigma

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// <----- Convert ----------------------------------------------------------------------------->
// <------------------------------------------------------------------------------------------->

// toMapAny converts any data into map[string]any.
func toMapAny(data any) map[string]any {
	newMap := map[string]any{}
	rtype := reflect.TypeOf(data)

	if rtype.Kind() == reflect.Map {
		val := reflect.ValueOf(data)

		for _, e := range val.MapKeys() {
			if k, ok := e.Interface().(string); ok {
				newMap[k] = val.MapIndex(e).Interface()
			}
		}
	}
	return newMap
}

// toSliceAny converts any data into []any.
func toSliceAny(data any) []any {
	newSlice := []any{}
	rtype := reflect.TypeOf(data)

	if rtype.Kind() == reflect.Slice {
		val := reflect.ValueOf(data)

		for i := 0; i < val.Len(); i++ {
			newSlice = append(newSlice, val.Index(i).Interface())
		}
	}
	return newSlice
}

// <----- Query Value ------------------------------------------------------------------------->
// <------------------------------------------------------------------------------------------->

// getValueMapPaths finds paths in nested data that end with a target value.
func getValueMapPaths(data map[string]any, pathSeen []string, target any, results *[][]string) {
	for k, v := range data {
		if reflect.DeepEqual(v, target) {
			nullPath := append(pathSeen, k)
			npCopy := make([]string, len(nullPath))
			copy(npCopy, nullPath)
			*results = append(*results, npCopy)
			continue
		}
		val := reflect.ValueOf(v)
		if val.Type().Kind() == reflect.Map {
			nullPath := append(pathSeen, k)
			getValueMapPaths(toMapAny(v), nullPath, target, results)
		}
		if val.Type().Kind() == reflect.Slice {
			nullPath := append(pathSeen, k)
			getValueSlicePaths(toSliceAny(v), nullPath, target, results)
		}
	}
}

// getValueSlicePaths finds paths in nested data that end with a target value.
func getValueSlicePaths(data []any, pathSeen []string, target any, results *[][]string) {
	for i := range data {
		if reflect.DeepEqual(data[i], target) {
			nullPath := append(pathSeen, fmt.Sprintf("%d", i))
			npCopy := make([]string, len(nullPath))
			copy(npCopy, nullPath)
			*results = append(*results, npCopy)
			continue
		}
		kind := reflect.TypeOf(data[i]).Kind()
		if kind == reflect.Map {
			nullPath := append(pathSeen, fmt.Sprintf("%d", i))
			getValueMapPaths(toMapAny(data[i]), nullPath, target, results)
		}
		if kind == reflect.Slice {
			nullPath := append(pathSeen, fmt.Sprintf("%d", i))
			getValueSlicePaths(toSliceAny(data[i]), nullPath, target, results)
		}
	}
}

// TODO traverse and record paths that end with the target value.
func getValueStructPaths(data any, pathSeen []string, target any, results *[][]string) {}

// <----- Query Path -------------------------------------------------------------------------->
// <------------------------------------------------------------------------------------------->

// TODO traverse data by path and return end value or an error if the path is invalid.
func getMapPathValue(data map[string]any, path []string) (any, error) {
	level := data
	for i, key := range path {
		if i == len(path)-1 {
			target, ok := level[key]
			if ok {
				return target, nil
			}
			return nil, errors.New("Invalid map path at: " + key)
		} else {
			nextLevel, ok := level[key]
			if !ok {
				return nil, errors.New("Invalid map path at: " + key)
			}
			val := reflect.ValueOf(nextLevel)
			if val.Type().Kind() == reflect.Slice {
				sl := toSliceAny(nextLevel)
				return getSlicePathValue(sl, path[i+1:])
			} else if val.Type().Kind() == reflect.Map {
				level = toMapAny(nextLevel)
			} else {
				return nil, errors.New("Invalid map value at: " + key)
			}
		}
	}
	return nil, errors.New("Invalid map path.")
}

// TODO traverse data by path and return end value or an error if the path is invalid.
func getSlicePathValue(data []any, path []string) (any, error) {
	level := data
	for i, key := range path {
		index, err := strconv.Atoi(key)
		if err != nil || index >= len(level) {
			return nil, errors.New("Invalid slice path at: " + key)
		}
		if i == len(path)-1 {
			return level[index], nil
		} else {
			nextLevel := level[index]
			val := reflect.ValueOf(nextLevel)
			if val.Type().Kind() == reflect.Slice {
				level = toSliceAny(nextLevel)
			} else if val.Type().Kind() == reflect.Map {
				return getMapPathValue(toMapAny(nextLevel), path[i+1:])
			} else {
				return nil, errors.New("Invalid slice value at: " + key)
			}
		}
	}
	return nil, errors.New("Invalid slice path.")
}

// TODO traverse data by path and return end value or an error if the path is invalid.
func getStructPathValue(data any, path []string) (any, error) {
	return nil, nil
}

// <----- Insert ------------------------------------------------------------------------------>
// <------------------------------------------------------------------------------------------->

// replaceMapValues traverses a deep map and replaces the end value with a new value.
func replaceMapValues(data *map[string]any, path []string, newValue any) error {
	level := data
	for i, key := range path {
		if i == len(path)-1 {
			(*level)[key] = newValue
		} else {
			nextLevel, ok := (*level)[key]
			if !ok {
				return errors.New("Invalid map path at: " + key)
			}
			val := reflect.ValueOf(nextLevel)
			if val.Type().Kind() == reflect.Slice {
				sl := toSliceAny(nextLevel)
				(*level)[key] = sl
				return replaceSliceValues(&sl, path[i+1:], newValue)
			} else if val.Type().Kind() == reflect.Map {
				nl := toMapAny(nextLevel)
				(*level)[key] = nl
				level = &nl
			} else {
				return errors.New("Invalid map value at: " + key)
			}
		}
	}
	return nil
}

// replaceSliceValues traverses a deep slice and replaces the end value with a new value.
func replaceSliceValues(data *[]any, path []string, newValue any) error {
	level := data
	for i, key := range path {
		index, err := strconv.Atoi(key)
		if err != nil || index >= len(*level) {
			return errors.New("Invalid slice path at: " + key)
		}
		if i == len(path)-1 {
			(*level)[index] = newValue
		} else {
			nextLevel := (*level)[index]
			val := reflect.ValueOf(nextLevel)
			if val.Type().Kind() == reflect.Slice {
				sl := toSliceAny(nextLevel)
				(*level)[index] = sl
				level = &sl
			} else if val.Type().Kind() == reflect.Map {
				nl := toMapAny(nextLevel)
				(*level)[index] = nl
				return replaceMapValues(&nl, path[i+1:], newValue)
			} else {
				return errors.New("Invalid slice value at: " + key)
			}
		}
	}
	return nil
}

// TODO traverses a deep struct and replaces the end value with a new value.
func replaceStructValues(data any, path []string, newValue any) error {
	return nil
}
