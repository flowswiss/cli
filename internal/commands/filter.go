package commands

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func stringValue(primitive reflect.Value) (string, error) {
	switch primitive.Kind() {
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		return strconv.FormatInt(primitive.Int(), 10), nil
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		return strconv.FormatUint(primitive.Uint(), 10), nil
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return strconv.FormatFloat(primitive.Float(), 'f', 5, 64), nil
	case reflect.String:
		return primitive.String(), nil
	default:
		return "", fmt.Errorf("unsupported type %s", primitive.Kind().String())
	}
}

func containsQuery(val reflect.Value, query string, depth, maxDepth int) bool {
	if depth >= maxDepth {
		return false
	}

	switch val.Kind() {
	case reflect.Ptr:
		return containsQuery(val.Elem(), query, depth, maxDepth)
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i)
			if containsQuery(item, query, depth+1, maxDepth) {
				return true
			}
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			item := val.MapIndex(key)
			if containsQuery(item, query, depth+1, maxDepth) {
				return true
			}
		}
	case reflect.Struct:
		if regexp.MustCompile("^\\d+$").MatchString(query) {
			idField := val.FieldByName("Id")
			if idField.IsValid() {
				return strconv.FormatUint(idField.Uint(), 10) == query
			}

			return false
		}

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if containsQuery(field, query, depth+1, maxDepth) {
				return true
			}
		}
	default:
		str, err := stringValue(val)
		if err != nil {
			return false
		}

		return strings.Contains(strings.ToLower(str), query)
	}

	return false
}

func filter(items interface{}, query string, maxDepth int) (interface{}, error) {
	query = strings.ToLower(query)

	val := reflect.ValueOf(items)
	if val.Kind() != reflect.Array && val.Kind() != reflect.Slice {
		return nil, fmt.Errorf("value must be of type array or slice, %s given", val.Kind().String())
	}

	result := reflect.New(val.Type()).Elem()

	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		if containsQuery(item, query, 0, maxDepth) {
			result = reflect.Append(result, item)
		}
	}

	return result.Interface(), nil
}

func findOne(items interface{}, query string, maxDepth int) (interface{}, error) {
	filtered, err := filter(items, query, maxDepth)
	if err != nil {
		return nil, err
	}

	val := reflect.ValueOf(filtered)

	if val.Len() == 0 {
		return nil, fmt.Errorf("no value found matching query %q", query)
	}

	if val.Len() > 1 {
		buf := &bytes.Buffer{}
		for i := 0; i < val.Len() && i < 3; i++ {
			buf.WriteString(fmt.Sprintf("%v", val.Index(i).Interface()))

			if i+1 < val.Len() {
				buf.WriteString(", ")
			}
		}

		if val.Len() > 3 {
			buf.WriteString("...")
		}

		return nil, fmt.Errorf("more than one matching result were found through query %q: %s", query, buf.String())
	}

	return val.Index(0).Interface(), nil
}
