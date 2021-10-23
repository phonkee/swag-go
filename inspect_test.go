package swag

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestInspectParams(t *testing.T) {
	t.Run("basic field validation", func(t *testing.T) {
		type Some struct {
			Param1 string
			Param2 int64 `json:"param2"`
			Param3 *int  `json:"param3"`
		}

		inspected := inspectParams(Some{}, spec.QueryParam)
		assert.NotNil(t, inspected)
		assert.Equal(t, inspected[0].In, "query")
		assert.True(t, len(inspected) == 3)
		assert.True(t, inspected[0].Name == "Param1")
		assert.True(t, inspected[0].Type == "string")

		assert.True(t, inspected[1].Name == "param2")
		assert.True(t, inspected[1].Type == "integer")

		assert.False(t, inspected[2].Required)

		// try path param
		inspected = inspectParams(Some{}, spec.PathParam)
		assert.NotNil(t, inspected)
		assert.Equal(t, inspected[0].In, "path")
	})

	t.Run("test required", func(t *testing.T) {
		type Some struct {
			NonRequired *int
		}
		inspected := inspectParams(Some{}, spec.QueryParam)
		assert.NotNil(t, inspected)
		assert.True(t, inspected[0].Required == false)
	})

	t.Run("sub structs naming with dots", func(t *testing.T) {
		type Second struct {
			Second1 string
			Second2 int `json:"second2"`
		}
		type Third struct {
			Third1 string
			Third2 int `json:"third2"`
		}
		type Fourth struct {
			Fourth1 string
			Fourth2 string `json:"fourth2"`
		}

		type Some struct {
			Second Second
			Third  Third `json:"third"`
			Fourth
		}

		inspected := inspectParams(Some{}, spec.PathParam)

		assert.Equal(t, "Second.Second1", inspected[0].Name)
		assert.Equal(t, "Second.second2", inspected[1].Name)
		assert.Equal(t, "third.Third1", inspected[2].Name)
		assert.Equal(t, "third.third2", inspected[3].Name)
		assert.Equal(t, "Fourth1", inspected[4].Name)
		assert.Equal(t, "fourth2", inspected[5].Name)
	})
}

func TestInspectSchema(t *testing.T) {
	t.Run("test invalid schema", func(t *testing.T) {
		invalid := []interface{}{
			1,
			true,
			nil,
		}
		for _, item := range invalid {
			defs := spec.Definitions{}
			assert.Panics(t, func() {
				inspectSchema(item, defs)
			})
		}
	})

	t.Run("test valid schema", func(t *testing.T) {
		type Response struct {
			Some []int
		}

		defs := spec.Definitions{}
		ref := inspectSchema(Response{}, defs)
		spew.Dump(ref)
	})

	t.Run("test top level array", func(t *testing.T) {
		type Response struct {
			Some []int
		}

		defs := spec.Definitions{}
		ref := inspectSchema([]Response{}, defs)
		_ = ref
	})

}
