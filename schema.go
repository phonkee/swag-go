package swag

import (
	"errors"
	"reflect"
	"sync"

	"github.com/go-openapi/spec"
)

var (
	schemaReg              = newSchemaRegistry()
	errSchemaValueNotFound = errors.New("not found schema")
)

func registerSchema(target interface{}, fn func(*schemaRegistry, interface{}) (*spec.Schema, error)) error {
	return schemaReg.registerSchema(target, fn)
}

func registerSchemaKind(target interface{}, fn func(*schemaRegistry, interface{}) (*spec.Schema, error)) error {
	return nil
}

func getSchema(target interface{}, defs spec.Definitions) (*spec.Schema, error) {
	return nil, nil
}

func newSchemaRegistry() *schemaRegistry {
	return &schemaRegistry{
		mutex:   sync.RWMutex{},
		storage: map[reflect.Type]func(*schemaRegistry, interface{}) (*spec.Schema, error){},
		kinds:   map[reflect.Kind]func(*schemaRegistry, interface{}) (*spec.Schema, error){},
	}
}

type schemaRegistry struct {
	mutex   sync.RWMutex
	storage map[reflect.Type]func(*schemaRegistry, interface{}) (*spec.Schema, error)
	kinds   map[reflect.Kind]func(*schemaRegistry, interface{}) (*spec.Schema, error)
}

func (s *schemaRegistry) registerSchema(target interface{}, fn func(*schemaRegistry, interface{}) (*spec.Schema, error)) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	typ := reflect.TypeOf(target)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	s.storage[typ] = fn
	return nil
}
func (s *schemaRegistry) registerKind(target interface{}, fn func(*schemaRegistry, interface{}) (*spec.Schema, error)) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	typ := reflect.TypeOf(target)

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	s.kinds[typ.Kind()] = fn

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
		result, err := found(s, target)
		if err != nil {
			return nil, err
		}
		result.Nullable = !required
		return result, nil
	}

	if found, ok := s.kinds[typ.Kind()]; ok {
		return found(s, target)
	}

	return nil, errSchemaValueNotFound
}
