package swag

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func specType(name string) spec.StringOrArray {
	return spec.StringOrArray{name}
}

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
			IntValue    int
			IntPtrValue *int       `json:"int_ptr_value"`
			Time        time.Time  `json:"tyme"`
			Time2       *time.Time `json:"tyme2"`
			Buul        bool       `json:"buul"`
			Buul2       *bool      `json:"buul2"`
			IntArray    []int      `json:"int_array"`
		}
		sch, err := getSchema(&ExampleSchema{}, make(spec.Definitions))
		assert.NoError(t, err)
		assert.NotNil(t, sch)
		// increment this when doing changes to ExampleSchema
		assert.Equal(t, len(sch.Properties), 7)

		assert.Equal(t, specType("integer"), sch.Properties["IntValue"].Type)
		assert.Equal(t, false, sch.Properties["IntValue"].Nullable)
		assert.Equal(t, specType("integer"), sch.Properties["int_ptr_value"].Type)
		// TODO: fix this
		assert.Equal(t, true, sch.Properties["int_ptr_value"].Nullable)

		assert.Equal(t, specType("integer"), sch.Properties["tyme"].Type)
		assert.Equal(t, "date-time", sch.Properties["tyme"].Format)
		assert.Equal(t, false, sch.Properties["tyme"].Nullable)

		assert.Equal(t, specType("integer"), sch.Properties["tyme2"].Type)
		assert.Equal(t, "date-time", sch.Properties["tyme2"].Format)
		assert.Equal(t, true, sch.Properties["tyme2"].Nullable)

		assert.Equal(t, specType("boolean"), sch.Properties["buul"].Type)
		assert.Equal(t, false, sch.Properties["buul"].Nullable)
		assert.Equal(t, specType("boolean"), sch.Properties["buul2"].Type)
		assert.Equal(t, true, sch.Properties["buul2"].Nullable)

		// IntArray
		assert.Equal(t, specType("array"), sch.Properties["int_array"].Type)
		assert.Equal(t, specType("integer"), sch.Properties["int_array"].Items.Schema.Type)

	})

}
