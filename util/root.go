package util

import (
	"context"
	"math/big"
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
