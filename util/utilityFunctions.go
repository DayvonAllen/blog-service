package util

import (
	"strings"
)

func GenerateKey(value string, query string) string {
	var key strings.Builder

	for _, v := range strings.Fields(value) {
		key.WriteString(v)
	}

	key.WriteString(":")
	key.WriteString(query)

	return key.String()
}