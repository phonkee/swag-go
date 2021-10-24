package swag

import (
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
			assert.NoError(t, reg.registerSchema(item.in, func(*schemaRegistry, interface{}) (*spec.Schema, error) {
				return item.schemaFunc(item.in), nil
			}))
			// check get
			sch, err := reg.getSchema(item.in, make(spec.Definitions))
			assert.NoError(t, err)
			assert.NotNil(t, sch)
			item.expect(t, sch)
		}

	})

	t.Run("test kind", func(t *testing.T) {
		reg := newSchemaRegistry()
		assert.NoError(t, reg.registerKind(int(0), func(*schemaRegistry, interface{}) (*spec.Schema, error) {
			return &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: []string{"integer"},
				},
			}, nil
		}))

		sch, err := reg.getSchema(int(0), make(spec.Definitions))
		assert.NoError(t, err)
		assert.NotNil(t, sch)

	})

}
