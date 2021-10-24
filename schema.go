package swag

import (
	"reflect"
	"sync"

	"github.com/go-openapi/spec"
)

var (
	schemaReg = newSchemaRegistry()
)

func registerSchema(target interface{}, fn func() (*spec.Schema, error)) error {
	return schemaReg.registerSchema(target, fn)
}

func getSchema(target interface{}, defs spec.Definitions) (*spec.Schema, error) {
	return nil, nil
}

func newSchemaRegistry() *schemaRegistry {
	return &schemaRegistry{
		mutex:   sync.RWMutex{},
		storage: map[reflect.Type]func() (*spec.Schema, error){},
	}
}

type schemaRegistry struct {
	mutex   sync.RWMutex
	storage map[reflect.Type]func() (*spec.Schema, error)
}

func (s *schemaRegistry) registerSchema(target interface{}, fn func() (*spec.Schema, error)) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	typ := reflect.TypeOf(target)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	s.storage[typ] = fn
	return nil
}

func (s *schemaRegistry) getSchema(target interface{}, defs spec.Definitions) (*spec.Schema, error) {
	panic("too not")
	return schemaReg.getSchema(target, defs)
}
