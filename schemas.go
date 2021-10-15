package swag

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/structtag"
	"github.com/go-openapi/spec"
)

func newSchemas() Schemas {
	return &schemas{
		storage: map[reflect.Type]*schemaInfo{},
	}
}

type schemaInfo struct {
	name   string
	ref    string
	null   bool
	schema *spec.Schema
}

type schemas struct {
	storage map[reflect.Type]*schemaInfo
	mutex   sync.RWMutex
}

func (s *schemas) MarshalJSON() ([]byte, error) {
	_ = spec.BooleanProperty()
	return []byte(""), nil
}

func (s *schemas) GetRef(i interface{}) string {
	typ := reflect.TypeOf(i)

	// no pointers
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// check if we have struct
	if typ.Kind() != reflect.Struct {
		panic(fmt.Sprintf("schema must be struct or pointer to struct, got %T: %v", i, i))
	}

	inspected := inspect(i)
	spew.Dump(inspected)

	return "/hello"
}

// inspect given value and returns appropriate schema
func inspect(i interface{}) *spec.Schema {
	typ := reflect.TypeOf(i)

	// something
	//s := spec.RefSchema()

	// no pointers
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tags, err := structtag.Parse(string(field.Tag))
		if err != nil {
			panic(err)
		}
		_ = tags
	}

	return nil
}
