package swag

import "github.com/go-openapi/spec"

type ResponseOptions struct {
	Description string `json:"description"`
}

func (r *ResponseOptions) Defaults() {

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

// inspectResponse inspects target and returns response definition
func inspectResponse(target interface{}) (*spec.Response, error) {
	return nil, nil
}
