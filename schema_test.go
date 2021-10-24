package swag

import (
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	t.Run("test basic types in schema", func(t *testing.T) {
		reg := newSchemaRegistry()
		assert.Nil(t, reg.registerSchema(int(1), func() (*spec.Schema, error) {
			return nil, nil
		}))

	})
}
