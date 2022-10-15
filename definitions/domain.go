package definitions

import "github.com/go-openapi/spec"

type Definitions interface {
	Register(what interface{}) spec.Schema
	Spec() spec.Definitions
}
