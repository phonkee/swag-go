package swag

import (
	"strings"

	"github.com/fatih/structs"
)

func getFieldName(field *structs.Field) string {
	result := field.Name()
	if tag := strings.TrimSpace(field.Tag("json")); tag != "" {
		return strings.Split(tag, ",")[0]
	}
	return result
}

func isFieldAvailable(field *structs.Field) bool {
	if tag := strings.TrimSpace(field.Tag("json")); tag != "" {
		for _, part := range strings.Split(tag, ",") {
			if part == "-" {
				return false
			}
		}
	}
	return true
}

func getFieldDescription(field *structs.Field) string {
	return strings.TrimSpace(field.Tag("swag_description"))
}
