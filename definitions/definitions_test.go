package definitions

import (
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

type someTestStruct struct {
	ID         *int
	NullableFK *int `json:"nullable_fk,omitempty"`
}

func TestDefinitions(t *testing.T) {
	t.Run("test basic types", func(t *testing.T) {
		d := New()
		assert.NotNil(t, d)

		data := []struct {
			in     interface{}
			typ    string
			format string
		}{
			{int(1), "numeric", "int64"}, // assume we are on 64 bit
			{int8(1), "numeric", "int8"},
			{int16(1), "numeric", "int16"},
			{int32(1), "numeric", "int32"},
			{int64(1), "numeric", "int64"},
			{"hello", "string", ""},
			{true, "boolean", ""},
			{false, "boolean", ""},
		}
		for _, item := range data {
			result := New().Register(item.in)
			assert.Equal(t, spec.StringOrArray{item.typ}, result.SchemaProps.Type)
			if result.Format != "" {
				assert.Equal(t, item.format, result.Format)
			}
		}
	})

	t.Run("test nullable", func(t *testing.T) {
		d := New()
		assert.NotNil(t, d)

		result := d.Register(&someTestStruct{})
		assert.Equal(t, "definitions.someTestStruct", result.Ref.String())
		assert.True(t, result.Nullable)
	})

}
