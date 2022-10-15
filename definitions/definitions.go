package definitions

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/fatih/structs"
	"github.com/go-openapi/spec"
	"github.com/phonkee/swag-go/utils"
)

// New instantiates new definitions
func New() Interface {
	result := &definitions{
		definitions: make(spec.Definitions),
		types:       make(map[reflect.Type]func(schema *spec.Schema)),
	}

	// provide custom default types
	result.RegisterType(reflect.TypeOf(time.Time{}), func(schema *spec.Schema) {
		schema.Type = []string{"integer"}
		schema.Format = "date-time"
	})

	return result
}

// definitions implements Interface interface
type definitions struct {
	definitions spec.Definitions
	types       map[reflect.Type]func(schema *spec.Schema)
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

	// first we check for custom type
	if fn, ok := d.types[typ]; ok {
		fn(result)
		return *result
	}

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
	case reflect.Array, reflect.Slice:
		result.Type = []string{"array"}
		tmp := d.Register(reflect.New(typ.Elem()).Elem().Interface())
		result.Items = &spec.SchemaOrArray{
			Schema: &tmp,
		}
	case reflect.Struct:
		// TODO: add support for embedded structs
		result.Type = []string{"object"}
		result.Properties = spec.SchemaProperties{}

		// now iterate over all fields, decide which needs to be added (json tags used)
		for _, field := range structs.New(what).Fields() {
			if !utils.IsFieldAvailable(field) {
				continue
			}
			if !field.IsExported() {
				continue
			}
			if field.IsEmbedded() {
				// here we need to call Register on this field, then iterate over all fields, add them to this struct,
				// and then delete this type from spec definitions (we don't want to pollute it)

				emb := d.Register(field.Value())
				for name, field := range emb.Properties {
					result.Properties[name] = field
				}
				delete(d.definitions, emb.ID)
				continue
			}
			name := utils.GetFieldName(field)

			// we need to handle pointer type here by ourselves, since reflect.New is strict about it
			ptr := false
			fieldType := reflect.TypeOf(field.Value())
			for {
				if fieldType.Kind() == reflect.Ptr {
					ptr = true
					fieldType = fieldType.Elem()
				} else {
					break
				}
			}

			fr := d.Register(reflect.New(fieldType).Interface())
			if ptr {
				fr.Nullable = true
			}
			result.Properties[name] = fr
		}
	}

	// assign to definitions
	d.definitions[id] = *result

	return *result
}

// RegisterType registers custom type that provides custom marshalling and unmarshalling
// Warning! This makes it higher priority than any other type.
// Warning! you cannot make different implementation for pointer. Pointer is handled by default.
func (d *definitions) RegisterType(what reflect.Type, fn func(schema *spec.Schema)) {
	for {
		if what.Kind() == reflect.Ptr {
			what = what.Elem()
		} else {
			break
		}
	}
	d.types[what] = fn
}

// Spec returns pointer to raw spec.Definitions
func (d *definitions) Spec() spec.Definitions {
	return d.definitions
}
