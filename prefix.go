package swag

import "github.com/go-openapi/spec"

type prefixInfo struct {
	swagger    Swagger
	pathPrefix string
}

func newPrefix(info *prefixInfo) Prefix {
	return &prefix{info: info}
}

// prefix is special path that only creates new paths
type prefix struct {
	info *prefixInfo
}

func (p *prefix) Path(path string, method string, options ...*PathOptions) Path {
	panic("implement me")
}

func (p *prefix) Response(status int, what ...interface{}) Path {
	panic("implement me")
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

func (p *prefix) Spec() spec.Paths {
	panic("implement me")
}
