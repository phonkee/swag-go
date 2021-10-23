package swag

import "github.com/go-openapi/spec"

// inspectParams inspects given target and calls fn callback on each parameter
func inspectParams(target interface{}, fn func(p *spec.Parameter)) []*spec.Parameter {
	return nil
}
