package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

// returns a pointer to the bool value
func Bool(v bool) *bool {
	return &v
}

// returns a pointer to the string value
func String(v string) *string {
	return &v
}

// returns a pointer to the string value. if v is empty, return nil
func StringIgnoreEmpty(v string) *string {
	if len(v) < 1 {
		return nil
	}
	return &v
}

// Int returns a pointer to the int value
func Int(v int) *int {
	return &v
}

// Int32 returns a pointer to the int32 value
func Int32(v int32) *int32 {
	return &v
}

// Int returns a pointer to the int value. if v is empty, return nil
func IntIgnoreEmpty(v int) *int {
	if v == 0 {
		return nil
	}
	return &v
}

// Int32 returns a pointer to the int32 value. if v is empty, return nil
func Int32IgnoreEmpty(v int32) *int32 {
	if v == 0 {
		return nil
	}
	return &v
}

// Int32 returns a pointer to the int32 value
func Int64IgnoreEmpty(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

// Float64 returns a pointer to the float64 value
func Float64(v float64) *float64 {
	return &v
}

// StringToInt convert the string to int, and return the pointer of int value
func StringToInt(i *string) *int {
	if i == nil || len(*i) == 0 {
		return nil
	}

	r, err := strconv.Atoi(*i)
	if err != nil {
		log.Printf("[ERROR] convert the string %q to int failed.", *i)
	}
	return &r
}

// StringToBool convert the string to boolean, and return the pointer of boolean value
func StringToBool(v interface{}) *bool {
	if v, ok := v.(string); ok {
		b, err := strconv.ParseBool(v)
		if err != nil {
			log.Printf("[ERROR] convert the string %q to boolean failed.", v)
		}

		return &b
	}

	return nil
}

// StringValue returns the string value
func StringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

// ValueIgnoreEmpty returns to the string value. if v is empty, return nil
func ValueIgnoreEmpty(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	vl := reflect.ValueOf(v)

	if !vl.IsValid() {
		log.Printf("[ERROR] The value (%#v) is invalid", v)
		return nil
	}

	if (vl.Kind() != reflect.Bool) && vl.IsZero() {
		return nil
	}

	if (vl.Kind() == reflect.Array || vl.Kind() == reflect.Slice) && vl.Len() == 0 {
		return nil
	}

	return v
}

// Try to parse the string value as the JSON format, if the operation failed, returns an empty map result.
func StringToJson(jsonStrObj string) interface{} {
	if jsonStrObj == "" {
		return nil
	}
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStrObj), &jsonMap)
	if err != nil {
		log.Printf("[ERROR] Unable to convert the JSON string to the map object: %s", err)
	}
	return jsonMap
}

// Try to convert the JSON object to the string value, if the operation failed, returns an empty string.
// And supports remove some special (nested) keys.
func JsonToString(jsonObj interface{}, keysToRemove ...string) string {
	if jsonObj == nil {
		return ""
	}

	var err error
	for _, keyToRemove := range keysToRemove {
		jsonObj, err = removeNestedKey(jsonObj, keyToRemove)
		if err != nil {
			log.Printf("[ERROR] Unable to remove the nested key from the JSON object: %s", err)
		}
	}
	jsonStr, err := json.Marshal(jsonObj)
	if err != nil {
		log.Printf("[ERROR] Unable to convert the JSON object to string: %s", err)
	}
	return string(jsonStr)
}

// removeNestedKey removes a nested key (also can be normal key) from a JSON object.
// It returns the modified map as an interface{}, or an error if the operation fails.
func removeNestedKey(jsonObj interface{}, nestedKey string) (interface{}, error) {
	keys := strings.Split(nestedKey, ".")
	if len(keys) == 0 {
		return jsonObj, fmt.Errorf("invalid nested key: %s", nestedKey)
	}

	// Type assert the input to a map, if possible.
	m, ok := jsonObj.(map[string]interface{})
	if !ok {
		return jsonObj, fmt.Errorf("the function input (jsonObj) is not a map[string]interface{}")
	}

	var remove func(map[string]interface{}, []string)
	remove = func(currentMap map[string]interface{}, keys []string) {
		// If there's only one key left, check if it exists and delete it.
		if len(keys) == 1 {
			_, exists := currentMap[keys[0]]
			if exists {
				delete(currentMap, keys[0])
			}
			return
		}

		// Look for the next key in the map.
		nextKey := keys[0]
		nextValue, exists := currentMap[nextKey]
		if !exists {
			return
		}

		nextMap, ok := nextValue.(map[string]interface{})
		if !ok {
			return
		}

		// Recurse into the next map with the remaining keys.
		remove(nextMap, keys[1:])
	}

	remove(m, keys)

	return m, nil
}
