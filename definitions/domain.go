package definitions

import (
	"reflect"

	"github.com/go-openapi/spec"
)

type Interface interface {
	Register(what interface{}) spec.Schema
	RegisterType(what reflect.Type, fn func(schema *spec.Schema))
	Spec() spec.Definitions
}
