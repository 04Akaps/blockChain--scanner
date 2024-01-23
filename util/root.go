package util

import (
	"context"
	"encoding/json"
	"math/big"
	"reflect"
	"strconv"
)

func Context() context.Context {
	return context.Background()
}

func ToString(req interface{}) string {
	switch req.(type) {
	case int:
		return strconv.Itoa(req.(int))
	case int64:
		return strconv.Itoa(int(req.(int64)))
	case *big.Int:
		return strconv.Itoa(int(req.(*big.Int).Int64()))
	case uint64:
		return strconv.Itoa(int(req.(uint64)))
	default:
		return ""
	}
}

func ToJson(t interface{}) (interface{}, error) {
	var v interface{}
	if bytes, err := json.Marshal(t); err != nil {
		return nil, err
	} else if err := json.Unmarshal(bytes, &v); err != nil {
		return nil, err
	} else {
		jsonMap := v.(map[string]interface{})
		for key, value := range jsonMap {
			if reflect.TypeOf(value) == reflect.TypeOf(float64(0)) {
				jsonMap[key] = int64(value.(float64))
			}
		}

		return jsonMap, nil
	}
}
