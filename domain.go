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

// Swagger is main interface
// it is returned from New call and means single service.
type Swagger interface {
	http.Handler
	json.Marshaler

	PathProvider

	// Prefix adds ability to group endpoints and have common properties (response, query params, path params)
	Prefix(path string) Prefix
}

type PathOptions struct {
	Description string
	ID          string
}

type Path interface {
	// Body is request body
	Body(interface{}) Path

	// PathParams adds path params
	PathParams(interface{}) Path

	// QueryParams params
	QueryParams(interface{}) Path

	// Response returned for given status code
	Response(status int, what ...interface{}) Path

	// Spec returns spec compatible Paths
	Spec() spec.Paths
}

// Prefix TODO: implement prefix in future
type Prefix interface {
	PathProvider

	// PathParams adds path params
	PathParams(interface{}) Path

	// QueryParams params
	QueryParams(interface{}) Path

	// Response returned for given status code
	Response(status int, what ...interface{}) Path
}
