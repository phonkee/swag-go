package swag

import "github.com/go-openapi/spec"

func newPath(p string, method string, c Components, options *PathOptions) *path {
	return &path{
		components: c,
		path:       p,
		method:     method,
		options:    options,
		responses:  map[int]*response{},
	}
}

type path struct {
	path       string
	method     string
	components Components
	options    *PathOptions
	responses  map[int]*response
}

type response struct {
	schema *spec.Schema
}

func (p *path) Params(i interface{}) Path {
	return p
}

func (p *path) Response(status int, what interface{}) Path {
	_ = p.components.GetSchema(what)



	return p
}
