package definitions

import (
	"reflect"

	"github.com/go-openapi/spec"
)

type Definitions interface {
	Register(what interface{}) spec.Schema
	RegisterType(what reflect.Type, fn func(schema *spec.Schema))
	UpdateSpec(*spec.Swagger) error
}
