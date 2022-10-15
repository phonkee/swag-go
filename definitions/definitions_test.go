package definitions

import (
	"reflect"
	"testing"
	"time"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestDefinitions(t *testing.T) {

	t.Run("test embedded struct", func(t *testing.T) {
		type Else struct {
			ID int `json:"id"`
		}
		type Something struct {
			Else
		}
		sch := New().Register(Something{})
		assert.Equal(t, spec.StringOrArray{"integer"}, sch.Properties["id"].Type)
	})

	t.Run("test struct", func(t *testing.T) {
		type _SomethingTesting struct {
			// test here
			ID          *int `json:"id,omitempty"`
			Ignored     *int `json:"-"`
			notExported *int
		}
		d := New()
		sch := d.Register(_SomethingTesting{})
		assert.Equal(t, spec.StringOrArray{"integer"}, sch.Properties["id"].Type)
		_, ok := sch.Properties["Ignored"]
		assert.False(t, ok)
	})

	t.Run("test unsupported type", func(t *testing.T) {
		assert.Panics(t, func() {
			New().Register(map[string]string{})
		})

	})

	t.Run("test custom type", func(t *testing.T) {
		type Custom struct {
		}

		for _, regType := range []any{Custom{}, &Custom{}} {
			d := New()
			d.RegisterType(reflect.TypeOf(regType), func(schema *spec.Schema) {
				schema.Type = []string{"custom"}
			})

			for _, item := range []struct {
				in       interface{}
				id       string
				nullable bool
			}{
				{Custom{}, "definitions.Custom", false},
				{&Custom{}, "definitions.Custom", true},
			} {
				regged := d.Register(item.in)
				assert.Equal(t, spec.StringOrArray{"custom"}, regged.Type)
				assert.Equal(t, item.id, regged.Ref.String())
				assert.Equal(t, item.nullable, regged.Nullable)
			}
		}

	})

	t.Run("test basic types", func(t *testing.T) {
		d := New()
		assert.NotNil(t, d)

		data := []struct {
			in       interface{}
			typ      string
			format   string
			itemsTyp string
		}{
			{in: int(1), typ: "integer", format: "int64"}, // assume we are on 64 bit
			{in: int8(1), typ: "integer", format: "int8"},
			{in: int16(1), typ: "integer", format: "int16"},
			{in: int32(1), typ: "integer", format: "int32"},
			{in: int64(1), typ: "integer", format: "int64"},
			{in: uint(1), typ: "integer", format: "int64"}, // assume we are on 64 bit
			{in: uint8(1), typ: "integer", format: "int8"},
			{in: uint16(1), typ: "integer", format: "int16"},
			{in: uint32(1), typ: "integer", format: "int32"},
			{in: uint64(1), typ: "integer", format: "int64"},
			{in: "hello", typ: "string"},
			{in: true, typ: "boolean"},
			{in: false, typ: "boolean"},
			{in: []int{1}, typ: "array", itemsTyp: "integer"},
			{in: []string{"hello"}, typ: "array", itemsTyp: "string"},
		}
		for _, item := range data {
			result := New().Register(item.in)
			assert.Equal(t, spec.StringOrArray{item.typ}, result.SchemaProps.Type)
			if result.Format != "" {
				assert.Equal(t, item.format, result.Format)
			}
			if item.itemsTyp != "" {
				assert.Equal(t, spec.StringOrArray{item.itemsTyp}, result.Items.Schema.SchemaProps.Type)
			}
		}
	})

	t.Run("test builtin custom types", func(t *testing.T) {
		data := []struct {
			in     interface{}
			typ    string
			format string
		}{
			{time.Time{}, "integer", "date-time"},
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
		type someTestStruct struct {
			ID         *int
			NullableFK *int `json:"nullable_fk,omitempty"`
		}

		d := New()
		assert.NotNil(t, d)

		result := d.Register(&someTestStruct{})
		assert.Equal(t, "definitions.someTestStruct", result.Ref.String())
		assert.True(t, result.Nullable)
	})

}
