package swag

import (
	"errors"
	"reflect"
	"sync"

	"github.com/go-openapi/spec"
)

var (
	globalParameters     = newParameters()
	errParameterNotFound = errors.New("parameter not found")
)

func init() {
	paramIntFunc := func(parameter *spec.Parameter, p reflect.Type) {
		parameter.Type = "integer"
	}
	registerParameter(paramIntFunc, reflect.TypeOf(int(0)), reflect.TypeOf(int8(0)), reflect.TypeOf(int16(0)), reflect.TypeOf(int32(0)), reflect.TypeOf(int64(0)))
	registerParameter(paramIntFunc, reflect.TypeOf(uint(0)), reflect.TypeOf(uint8(0)), reflect.TypeOf(uint16(0)), reflect.TypeOf(uint32(0)), reflect.TypeOf(uint64(0)))
	registerParameter(func(parameter *spec.Parameter, t reflect.Type) {
		switch t.Kind() {
		case reflect.Float32:
			parameter.Type = "number"
			parameter.Format = "float"
		case reflect.Float64:
			parameter.Type = "number"
			parameter.Format = "double"
		}
	}, reflect.TypeOf(float32(1)), reflect.TypeOf(float64(1)))
	registerParameter(func(parameter *spec.Parameter, r reflect.Type) {
		parameter.Type = "string"
	}, reflect.TypeOf(""))
}

// registerParameter registers parameter for given type
func registerParameter(fn func(*spec.Parameter, reflect.Type), types ...reflect.Type) {
	globalParameters.RegisterParameter(fn, types...)
}

// getParameter returns parameter
func getParameter(typ reflect.Type, parameter *spec.Parameter) (*spec.Parameter, error) {
	return globalParameters.Get(typ, parameter)
}

// newParameters storage
func newParameters() *parameters {
	return &parameters{
		mutex:   sync.RWMutex{},
		storage: map[reflect.Type]func(*spec.Parameter, reflect.Type){},
	}
}

type parameters struct {
	mutex   sync.RWMutex
	storage map[reflect.Type]func(*spec.Parameter, reflect.Type)
}

func (p *parameters) Get(typ reflect.Type, param *spec.Parameter) (*spec.Parameter, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	required := true

	// pointer means non required, so we just set Required property
	for typ.Kind() == reflect.Ptr {
		required = false
		typ = typ.Elem()
	}

	if fn, ok := p.storage[typ]; ok {
		param.Required = required
		fn(param, typ)
		return param, nil
	}

	return nil, errParameterNotFound
}

func (p *parameters) RegisterParameter(fn func(*spec.Parameter, reflect.Type), types ...reflect.Type) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, t := range types {
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		p.storage[t] = fn
	}
}
