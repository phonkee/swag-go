package swag

import (
	"reflect"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func intPtr(i int) *int {
	return &i
}

func TestParameters(t *testing.T) {
	t.Run("test valid parameters", func(t *testing.T) {
		data := []struct {
			typ        reflect.Type
			fn         func(parameter *spec.Parameter, p reflect.Type)
			example    interface{}
			expectType string
		}{
			{reflect.TypeOf(int(10)), func(parameter *spec.Parameter, _ reflect.Type) { parameter.Type = "integer" }, int(1), "integer"},
			{reflect.TypeOf(int(10)), func(parameter *spec.Parameter, _ reflect.Type) { parameter.Type = "integer" }, intPtr(1), "integer"},
		}

		for _, item := range data {
			ps := newParameters()
			ps.RegisterParameter(item.fn, item.typ)
			p, err := ps.Get(reflect.TypeOf(item.example), spec.QueryParam("param"))
			assert.NoError(t, err)
			assert.Equal(t, p.Type, item.expectType)
		}
	})

	t.Run("test custom parameters", func(t *testing.T) {
		type Example struct {
		}
		ps := newParameters()
		ps.RegisterParameter(func(parameter *spec.Parameter, r reflect.Type) {
			parameter.Type = "example"
			parameter.Format = "custom"
		}, reflect.TypeOf(Example{}))

		ex, err := ps.Get(reflect.TypeOf(&Example{}), spec.QueryParam("example"))
		assert.NoError(t, err)
		assert.Equal(t, ex.Type, "example")
		assert.Equal(t, ex.Format, "custom")
		assert.Equal(t, ex.Required, false)
	})
}
