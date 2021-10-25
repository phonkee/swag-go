package swag

// TODO: prefix is not implemented yet, should not be very complicated though
// Currently it's lower priority but soon it will be implemented

type PrefixOptions struct {
}

func (p *PrefixOptions) Defaults() {

}

type prefixInfo struct {
	swagger    *swagger
	pathPrefix string
	resetCache func()
}

func newPrefix(info *prefixInfo, options *PrefixOptions) Prefix {
	info.resetCache()
	if options == nil {
		options = &PrefixOptions{}
	}
	return &prefix{info: info, options: options}
}

// prefix is special path that only creates new paths
type prefix struct {
	info    *prefixInfo
	options *PrefixOptions
}

func (p *prefix) Prefix(path string, options ...*PrefixOptions) Prefix {
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
