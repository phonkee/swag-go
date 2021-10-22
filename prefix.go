package swag

import "github.com/go-openapi/spec"

// prefix is special path that only creates new paths
type prefix struct {
}

func (p *prefix) Body(i interface{}) Path {
	panic("implement me")
}

func (p *prefix) PathParams(i interface{}) Path {
	panic("implement me")
}

func (p *prefix) QueryParams(i interface{}) Path {
	panic("implement me")
}

func (p *prefix) Response(status int, what interface{}) Path {
	panic("implement me")
}

func (p *prefix) Spec() spec.Paths {
	panic("implement me")
}
