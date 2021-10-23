package swag

import "github.com/go-openapi/spec"

type prefixInfo struct {
	swagger    *swagger
	pathPrefix string
	resetCache func()
}

func newPrefix(info *prefixInfo) Prefix {
	info.resetCache()
	return &prefix{info: info}
}

// prefix is special path that only creates new paths
type prefix struct {
	info *prefixInfo
}

func (p *prefix) Prefix(path string) Prefix {
	p.info.resetCache()
	panic("implement me")
}

func (p *prefix) Path(path string, method string, options ...*PathOptions) Path {
	p.info.resetCache()
	panic("implement me")
}

func (p *prefix) Response(status int, response interface{}, options ...*ResponseOptions) Path {
	p.info.resetCache()
	panic("implement me")
}

func (p *prefix) Body(i interface{}) Path {
	p.info.resetCache()
	panic("implement me")
}

func (p *prefix) PathParams(i interface{}) Path {
	p.info.resetCache()
	panic("implement me")
}

func (p *prefix) QueryParams(i interface{}) Path {
	p.info.resetCache()
	panic("implement me")
}

func (p *prefix) Spec() spec.Paths {
	panic("implement me")
}
