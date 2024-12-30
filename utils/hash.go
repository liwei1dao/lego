package utils

import (
	"fmt"
	"hash/fnv"
)

func HashString(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func GenHashKey(options ...interface{}) uint64 {
	keyString := ""
	for _, opt := range options {
		keyString = keyString + "/" + toString(opt)
	}

	return HashString(keyString)
}

func toString(obj interface{}) string {
	return fmt.Sprintf("%v", obj)
}
