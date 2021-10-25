package swag

import (
	"strings"

	"github.com/go-openapi/spec"
)

type ResponseOptions struct {
	Description string
	// Headers is any structure with fields
	Headers interface{}
}

func (r *ResponseOptions) Defaults() {
	r.Description = strings.TrimSpace(r.Description)
}

type response struct {
	options *ResponseOptions
	status  int
	target  interface{}
}

func newResponse(status int, target interface{}, options *ResponseOptions) *response {
	if options == nil {
		options = &ResponseOptions{}
	}
	// set defaults
	options.Defaults()

	// return response
	return &response{
		options: options,
		status:  status,
		target:  target,
	}
}

func (r *response) spec() *spec.Response {
	return &spec.Response{}
}
