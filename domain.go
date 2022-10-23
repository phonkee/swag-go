package swag

import (
	"encoding/json"
	"net/http"

	"github.com/go-openapi/spec"
)

type PathProvider interface {
	// Path adds new endpoints
	Path(path string, method string, options ...*PathOptions) Path
}

type PrefixProvider interface {
	// Prefix adds ability to group endpoints and have common properties (response, query pathParams, pathImpl pathParams)
	Prefix(path string, options ...*PrefixOptions) Prefix
}

// Swag is main interface
// it is returned from New call and means single service.
type Swag interface {
	http.Handler
	json.Marshaler

	PathProvider
	PrefixProvider

	Debug()

	// RegisterType registers type that has special marshaling
	RegisterType(what interface{}, fn func(schema *spec.Schema))

	// Response returned for given status code
	// if no response is provided, no body is defined, if only nil is passed all previous responses defined will be removed)
	Response(status int, response interface{}, options ...*ResponseOptions) Swag

	// spec returns specification and caches it
	spec() *spec.Swagger
}

type Path interface {
	UpdateSpec
	// Body is request body
	Body(interface{}) Path

	// PathParams adds pathImpl pathParams (if nil is provided, all pathParams previously defined will be removed)
	PathParams(interface{}) Path

	// QueryParams pathParams (if nil is provided, all pathParams previously defined will be removed)
	QueryParams(interface{}) Path

	// Response returned for given status code
	// if no response is provided, no body is defined, if only nil is passed all previous responses defined will be removed)
	Response(status int, response interface{}, options ...*ResponseOptions) Path
}

// Prefix TODO: implement prefix in future
type Prefix interface {
	PathProvider
	PrefixProvider
	UpdateSpec

	// PathParams adds pathImpl pathParams
	PathParams(interface{}) Prefix

	// QueryParams pathParams
	QueryParams(interface{}) Prefix

	// Response returned for given status code
	Response(status int, response interface{}, options ...*ResponseOptions) Prefix
}

// UpdateSpec interface is used to update specification
type UpdateSpec interface {
	UpdateSpec(*spec.Swagger) error
}
