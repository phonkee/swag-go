package swag

// TODO: prefix is not implemented yet, should not be very complicated though
// Currently it's lower priority but soon it will be implemented

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

func (p *prefix) Response(status int, response interface{}, options ...*ResponseOptions) Prefix {
	p.info.resetCache()
	panic("implement me")
}

func (p *prefix) PathParams(i interface{}) Prefix {
	p.info.resetCache()
	panic("implement me")
}

func (p *prefix) QueryParams(i interface{}) Prefix {
	p.info.resetCache()
	panic("implement me")
}
