package swag

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/fatih/structs"
	"github.com/go-openapi/spec"
)

var (
	schemaReg              = newSchemaRegistry()
	errSchemaValueNotFound = errors.New("not found schema")
)

func init() {
	// register pointer kind
	mustRegisterSchemaKind(func(registry *schemaRegistry, i interface{}, d spec.Definitions) (sch *spec.Schema, err error) {
		typ := reflect.TypeOf(i)
		for typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		return registry.getSchema(reflect.New(typ), d)
	}, reflect.Ptr)

	// register integer kinds
	intSchemaKindFunc := func(registry *schemaRegistry, i interface{}, d spec.Definitions) (sch *spec.Schema, err error) {
		return &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"integer"},
			},
		}, nil
	}
	mustRegisterSchemaKind(intSchemaKindFunc, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64)
	mustRegisterSchemaKind(intSchemaKindFunc, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64)

	// register struct kind
	mustRegisterSchemaKind(func(registry *schemaRegistry, i interface{}, definitions spec.Definitions) (*spec.Schema, error) {
		id := fmt.Sprintf("%T", i)
		result := spec.RefSchema(id)
		result.ID = id
		result.Properties = spec.SchemaProperties{}

		ss := structs.New(i)
		for _, field := range ss.Fields() {
			_ = field
		}
		return result, nil
	}, reflect.Struct)

}

func mustRegisterSchema(target interface{}, fn func(*schemaRegistry, interface{}, spec.Definitions) (*spec.Schema, error)) error {
	return schemaReg.registerSchema(target, fn)
}

func mustRegisterSchemaKind(fn func(*schemaRegistry, interface{}, spec.Definitions) (*spec.Schema, error), targets ...reflect.Kind) {
	if err := schemaReg.registerKind(fn, targets...); err != nil {
		panic(err)
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
func (s *schemaRegistry) registerKind(fn func(*schemaRegistry, interface{}, spec.Definitions) (*spec.Schema, error), targets ...reflect.Kind) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, kind := range targets {
		s.kinds[kind] = fn
	}

	return nil
}

func (s *schemaRegistry) getSchema(target interface{}, defs spec.Definitions) (*spec.Schema, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	typ := reflect.TypeOf(target)
	required := true
	for typ.Kind() == reflect.Ptr {
		required = false
		typ = typ.Elem()
	}

	if found, ok := s.storage[typ]; ok {
		result, err := found(s, target, defs)
		if err != nil {
			return nil, err
		}
		result.Nullable = !required
		return result, nil
	}

	if found, ok := s.kinds[typ.Kind()]; ok {
		return found(s, target, defs)
	}

	return nil, errSchemaValueNotFound
}
