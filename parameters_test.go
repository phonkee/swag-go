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
			kind       reflect.Kind
			fn         func(p *parameters, parameter *spec.Parameter, _ interface{})
			example    interface{}
			expectType string
		}{
			{
				reflect.TypeOf(int(10)).Kind(),
				func(p *parameters, parameter *spec.Parameter, _ interface{}) { parameter.Type = "integer" },
				int(1),
				"integer",
			},
		}

		for _, item := range data {
			ps := newParameters()
			assert.NoError(t, ps.RegisterKind([]reflect.Kind{item.kind}, item.fn))
			p, err := ps.Get(item.example, spec.QueryParam("param"))
			assert.NoError(t, err, "example: %v", item.example)
			assert.Equal(t, p.Type, item.expectType)
		}
	})

	t.Run("test custom parameters", func(t *testing.T) {
		type Example struct {
		}
		ps := newParameters()
		ps.RegisterParameter([]reflect.Type{reflect.TypeOf(Example{})}, func(p *parameters, parameter *spec.Parameter, r reflect.Type) {
			parameter.Type = "example"
			parameter.Format = "custom"
		})

		ex, err := ps.Get(Example{}, spec.QueryParam("example"))
		assert.NoError(t, err)
		assert.Equal(t, ex.Type, "example")
		assert.Equal(t, ex.Format, "custom")
		assert.Equal(t, ex.Required, false)
	})

	t.Run("test all registered", func(t *testing.T) {
		type FullParametersTestStruct struct {
			SomInt int
		}

		x := inspectParams(&FullParametersTestStruct{}, spec.QueryParam)
		assert.Equal(t, 1, len(x))
	})

}
