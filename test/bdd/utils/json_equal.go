package utils

import (
	"encoding/json"
	"reflect"
)

func JsonEqual(a, b string) (bool, error) {
	var j1, j2 map[string]interface{}

	if err := json.Unmarshal([]byte(a), &j1); err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(b), &j2); err != nil {
		return false, err
	}

	return reflect.DeepEqual(j1, j2), nil
}
