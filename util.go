package db

import (
	"fmt"
	"strings"
)

func inArray(val string, array []string) bool {
	for _, item := range array {
		if val == item {
			return true
		}
	}
	return false
}

func AddSlashes(str string, charlist string) string {
	for _, char := range charlist {
		str = strings.Replace(str, string(char), fmt.Sprintf("\\%s", string(char)), 0)
	}
	return str
}
