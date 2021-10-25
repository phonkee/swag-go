package swag

import (
	"fmt"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	t.Run("test basic types in schema", func(t *testing.T) {
		data := []struct {
			in         interface{}
			schemaFunc func(interface{}) *spec.Schema
			expect     func(t2 *testing.T, schema *spec.Schema)
		}{
			{int(1), func(interface{}) *spec.Schema {
				return &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: []string{"integer"},
					},
				}
			}, func(t *testing.T, s *spec.Schema) {

				// test here

			}},
		}

		for _, item := range data {
			reg := newSchemaRegistry()
			assert.NoError(t, reg.registerSchema(item.in, func(*schemaRegistry, interface{}, spec.Definitions) (*spec.Schema, error) {
				return item.schemaFunc(item.in), nil
			}))
			// check get
			sch, err := reg.getSchema(item.in, make(spec.Definitions))
			assert.NoError(t, err)
			assert.NotNil(t, sch)
			item.expect(t, sch)
		}

	})

	t.Run("test custom kind", func(t *testing.T) {
		reg := newSchemaRegistry()
		type _SomethingTesting struct{}

		assert.NoError(t, reg.registerSchema(_SomethingTesting{}, func(_ *schemaRegistry, in interface{}, d spec.Definitions) (*spec.Schema, error) {
			return spec.RefSchema(fmt.Sprintf("%T", in)), nil
		}))

		sch, err := reg.getSchema(_SomethingTesting{}, make(spec.Definitions))
		assert.NoError(t, err)
		assert.NotNil(t, sch)
		assert.Equal(t, "swag._SomethingTesting", sch.Ref.String())
	})

	t.Run("test defined kinds", func(t *testing.T) {
		type ExampleSchema struct {
			IntValue int
		}
		sch, err := getSchema(&ExampleSchema{}, make(spec.Definitions))
		assert.NoError(t, err)
		assert.NotNil(t, sch)
	})

}
