package definitions

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-openapi/spec"
)

// New instantiates new definitions
func New() Definitions {
	return &definitions{
		definitions: make(spec.Definitions),
	}
}

// definitions implements Definitions interface
type definitions struct {
	definitions spec.Definitions
}

// Register registers given type with given schema, this is called recursively and when definition is already registered
// it uses it. Due to the fact that this is recursive, we cannot protect map access to map with mutex, so it's not
// currently safe to use concurrently.
func (d *definitions) Register(what interface{}) spec.Schema {
	typ := reflect.TypeOf(what)
	val := reflect.ValueOf(what)

	// handle pointers here
	if typ.Kind() == reflect.Ptr {
		result := d.Register(val.Elem().Interface())
		result.Nullable = true
		return result
	}

	// prepare ref schema, first id and then assign
	id := fmt.Sprintf("%T", what)
	result := spec.RefSchema(id)
	result.ID = id
	result.Properties = spec.SchemaProperties{}

	// now check for type
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result.Type = []string{"integer"}
		switch typ.Kind() {
		case reflect.Int:
			result.Format = "int" + strconv.FormatInt(int64(typ.Bits()), 10)
		case reflect.Int8:
			result.Format = "int8"
		case reflect.Int16:
			result.Format = "int16"
		case reflect.Int32:
			result.Format = "int32"
		case reflect.Int64:
			result.Format = "int64"
		default:
			panic(fmt.Sprintf("don't know how to handle this type: %v", what))
		}
	case reflect.String:
		result.Type = []string{"string"}
	case reflect.Bool:
		result.Type = []string{"boolean"}
	}

	// assign to definitions
	d.definitions[id] = *result

	return *result
}

// Spec returns pointer to raw spec.Definitions
func (d *definitions) Spec() spec.Definitions {
	return d.definitions
}
