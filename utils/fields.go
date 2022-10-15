package utils

import (
	"strings"

	"github.com/fatih/structs"
)

func GetFieldName(field *structs.Field) string {
	result := field.Name()
	if tag := strings.TrimSpace(field.Tag("json")); tag != "" {
		return strings.Split(tag, ",")[0]
	}
	return result
}
