package swag

import "strings"

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
	options.Defaults()
	return &response{
		options: options,
		status:  status,
		target:  target,
	}
}
