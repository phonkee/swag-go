package swag

import (
	"github.com/go-openapi/spec"
)

// Password type
type Password string

func init() {
	// add support for password
	mustRegisterSchemaType(Password(""), func(registry *schemaRegistry, i interface{}, definitions spec.Definitions) (*spec.Schema, error) {
		return &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:   []string{"string"},
				Format: "password",
			},
		}, nil
	})
}
