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
	paramIntFunc := func(p *parameters, parameter *spec.Parameter, t reflect.Type) {
		parameter.Type = "integer"
	}
	registerParameterType([]reflect.Type{
		reflect.TypeOf(int(0)),
		reflect.TypeOf(int8(0)),
		reflect.TypeOf(int16(0)),
		reflect.TypeOf(int32(0)),
		reflect.TypeOf(int64(0)),
	}, paramIntFunc)
	registerParameterType([]reflect.Type{
		reflect.TypeOf(uint(0)),
		reflect.TypeOf(uint8(0)),
		reflect.TypeOf(uint16(0)),
		reflect.TypeOf(uint32(0)),
		reflect.TypeOf(uint64(0)),
	}, paramIntFunc)
	registerParameterType([]reflect.Type{
		reflect.TypeOf(float32(1)),
		reflect.TypeOf(float64(1)),
	}, func(p *parameters, parameter *spec.Parameter, t reflect.Type) {
		switch t.Kind() {
		case reflect.Float32:
			parameter.Type = "number"
			parameter.Format = "float"
		case reflect.Float64:
			parameter.Type = "number"
			parameter.Format = "double"
		}
	})
	registerParameterType([]reflect.Type{reflect.TypeOf("")},
		func(p *parameters, parameter *spec.Parameter, r reflect.Type) {
			parameter.Type = "string"
		},
	)

	registerParameterType([]reflect.Type{reflect.TypeOf(true)},
		func(p *parameters, parameter *spec.Parameter, r reflect.Type) {
			parameter.Type = "boolean"
		},
	)

	// TODO: finish Ptr
	registerParameterKind([]reflect.Kind{reflect.Ptr}, func(p *parameters, parameter *spec.Parameter, i interface{}) {
		typ := reflect.TypeOf(i)
		for typ.Kind() != reflect.Ptr {
			typ = typ.Elem()
		}

		panic("not implemented")
	})

	// TODO: finish struct and remove hardcoded
	//registerParameterKind([]reflect.Kind{reflect.Struct}, func(p *parameters, parameter *spec.Parameter, i interface{}) {
	//	panic("not implemented")
	//})

}

// registerParameterType registers parameter for given type
func registerParameterType(types []reflect.Type, fn func(*parameters, *spec.Parameter, reflect.Type)) {
	globalParameters.RegisterParameter(types, fn)
}

func registerParameterKind(kinds []reflect.Kind, fn func(*parameters, *spec.Parameter, interface{})) {
	if err := globalParameters.RegisterKind(kinds, fn); err != nil {
		panic(err)
	}
}

// getParameter returns parameter
func getParameter(typ reflect.Type, parameter *spec.Parameter) (*spec.Parameter, error) {
	return globalParameters.Get(typ, parameter)
}

// newParameters types
func newParameters() *parameters {
	return &parameters{
		mutex: sync.RWMutex{},
		types: map[reflect.Type]func(*parameters, *spec.Parameter, reflect.Type){},
		kinds: map[reflect.Kind]func(*parameters, *spec.Parameter, interface{}){},
	}
}

type parameters struct {
	mutex sync.RWMutex
	types map[reflect.Type]func(*parameters, *spec.Parameter, reflect.Type)
	kinds map[reflect.Kind]func(*parameters, *spec.Parameter, interface{})
}

func (p *parameters) Get(what interface{}, param *spec.Parameter) (*spec.Parameter, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	typ := reflect.TypeOf(what)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if fn, ok := p.types[typ]; ok {
		fn(p, param, typ)
		return param, nil
	}

	if fn, ok := p.kinds[typ.Kind()]; ok {
		fn(p, param, what)
		return param, nil
	}

	return nil, errParameterNotFound
}

func (p *parameters) RegisterParameter(types []reflect.Type, fn func(*parameters, *spec.Parameter, reflect.Type)) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, t := range types {
		p.types[t] = fn
	}
}

func (p *parameters) RegisterKind(kinds []reflect.Kind, fn func(*parameters, *spec.Parameter, interface{})) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, kind := range kinds {
		p.kinds[kind] = fn
	}

	return nil
}
