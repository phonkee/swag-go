package swag

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/fatih/structs"
	"github.com/go-openapi/spec"
)

var (
	schemaReg              = newSchemaRegistry()
	errSchemaValueNotFound = errors.New("not found schema")
)

func init() {
	// register pointer kind
	mustRegisterSchemaKind([]reflect.Kind{reflect.Ptr}, func(registry *schemaRegistry, i interface{}, d spec.Definitions) (sch *spec.Schema, err error) {
		typ := reflect.TypeOf(i)
		for typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		sch, err = registry.getSchema(reflect.New(typ).Elem().Interface(), d)
		if err != nil {
			return
		}

		sch.Nullable = true

		return sch, nil
	})

	// register integer kinds
	mustRegisterSchemaKind([]reflect.Kind{reflect.String}, func(registry *schemaRegistry, i interface{}, definitions spec.Definitions) (*spec.Schema, error) {
		return &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"string"},
			},
		}, nil
	})
	// register integer kinds
	intSchemaKindFunc := func(registry *schemaRegistry, i interface{}, d spec.Definitions) (sch *spec.Schema, err error) {
		return &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"integer"},
			},
		}, nil
	}
	mustRegisterSchemaKind([]reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64}, intSchemaKindFunc)
	mustRegisterSchemaKind([]reflect.Kind{reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64}, intSchemaKindFunc)

	// register boolean type
	mustRegisterSchemaKind([]reflect.Kind{reflect.Bool}, func(registry *schemaRegistry, i interface{}, definitions spec.Definitions) (*spec.Schema, error) {
		return &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"boolean"},
			},
		}, nil
	})

	// register struct kind
	mustRegisterSchemaKind([]reflect.Kind{reflect.Struct}, func(registry *schemaRegistry, i interface{}, definitions spec.Definitions) (*spec.Schema, error) {
		id := fmt.Sprintf("%T", i)
		result := spec.RefSchema(id)
		result.ID = id
		result.Properties = spec.SchemaProperties{}
		ss := structs.New(i)

		typ := reflect.TypeOf(i)
		for typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		for index, field := range ss.Fields() {
			if !isFieldAvailable(field) {
				continue
			}

			// TODO: pointers not working?
			name := getFieldName(field)

			sch, err := registry.getSchema(field.Value(), definitions)
			if err != nil {
				return nil, err
			}

			sch.ID = name

			// TODO: hack for now
			if typ.Field(index).Type.Kind() == reflect.Ptr {
				sch.Nullable = true
			}

			result.Properties[name] = *sch
		}
		return result, nil
	})

	// handle time.Time
	mustRegisterSchemaType(time.Time{}, func(registry *schemaRegistry, i interface{}, definitions spec.Definitions) (*spec.Schema, error) {
		return &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:   []string{"integer"},
				Format: "date-time",
			},
		}, nil
	})

	// handle array/slice
	mustRegisterSchemaKind(
		[]reflect.Kind{reflect.Array, reflect.Slice},
		func(registry *schemaRegistry, i interface{}, definitions spec.Definitions) (*spec.Schema, error) {
			sch := &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: []string{"array"},
				},
			}

			elem := reflect.TypeOf(i).Elem()
			if elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}

			inner, errInner := registry.getSchema(reflect.New(elem).Elem().Interface(), definitions)
			if errInner != nil {
				return nil, errInner
			}

			sch.Items = &spec.SchemaOrArray{
				Schema: inner,
			}

			// TODO: remove this hack
			if reflect.TypeOf(i).Elem().Kind() == reflect.Ptr {
				sch.Items.Schema.Nullable = true
			}

			return sch, nil
		})
}

func mustRegisterSchemaType(target interface{}, fn func(*schemaRegistry, interface{}, spec.Definitions) (*spec.Schema, error)) {
	if err := schemaReg.registerSchema(target, fn); err != nil {
		panic(fmt.Sprintf("cannot register schema type %T: %v", target, err))
	}
}

func mustRegisterSchemaKind(kinds []reflect.Kind, fn func(*schemaRegistry, interface{}, spec.Definitions) (*spec.Schema, error)) {
	if err := schemaReg.registerKind(kinds, fn); err != nil {
		panic(fmt.Sprintf("cannot register schema kind: %v", err))
	}
}

func getSchema(target interface{}, defs spec.Definitions) (*spec.Schema, error) {
	return schemaReg.getSchema(target, defs)
}

func newSchemaRegistry() *schemaRegistry {
	return &schemaRegistry{
		mutex:   sync.RWMutex{},
		storage: map[reflect.Type]func(*schemaRegistry, interface{}, spec.Definitions) (*spec.Schema, error){},
		kinds:   map[reflect.Kind]func(*schemaRegistry, interface{}, spec.Definitions) (*spec.Schema, error){},
	}
}

type schemaRegistry struct {
	mutex   sync.RWMutex
	storage map[reflect.Type]func(*schemaRegistry, interface{}, spec.Definitions) (*spec.Schema, error)
	kinds   map[reflect.Kind]func(*schemaRegistry, interface{}, spec.Definitions) (*spec.Schema, error)
}

func (s *schemaRegistry) registerSchema(target interface{}, fn func(*schemaRegistry, interface{}, spec.Definitions) (*spec.Schema, error)) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	typ := reflect.TypeOf(target)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	s.storage[typ] = fn
	return nil
}
func (s *schemaRegistry) registerKind(kinds []reflect.Kind, fn func(*schemaRegistry, interface{}, spec.Definitions) (*spec.Schema, error)) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, kind := range kinds {
		s.kinds[kind] = fn
	}

	return nil
}

func (s *schemaRegistry) getSchema(target interface{}, defs spec.Definitions) (*spec.Schema, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	typ := reflect.TypeOf(target)

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if found, ok := s.storage[typ]; ok {
		result, err := found(s, target, defs)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	if found, ok := s.kinds[typ.Kind()]; ok {
		return found(s, target, defs)
	}

	return nil, errSchemaValueNotFound
}
