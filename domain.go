package swag

import "encoding/json"

type Service interface {
	json.Marshaler

	// Path adds new endpoints
	Path(path string, method string, options ...*PathOptions) Path
}

type ServiceOptions struct {
	Description string
}

type Path interface {
	// Body is request body
	Body(interface{}) Path

	// Params adds path params
	Params(interface{}) Path

	// Response returned for given status code
	Response(status int, what interface{}) Path
}

type PathOptions struct {
	Description string
}

type Components interface {
	json.Marshaler

	GetSchema(interface{}) string
}

type Schemas interface {
	json.Marshaler
	// Get stores schema and returns string to be returned
	GetRef(interface{}) string
}
