package swag

import (
	"strings"

	"github.com/go-openapi/spec"
)

// TODO: prefix is not implemented yet, should not be very complicated though
// Currently it's lower priority but soon it will be implemented

type PrefixOptions struct {
}

func (p *PrefixOptions) Defaults() {
}

type prefixInfo struct {
	definitions spec.Definitions
	swagger     *swagger
	pathPrefix  string
	resetCache  func()
	responses   map[int]*response
	invalidate  func()
}

func newPrefix(info *prefixInfo, options *PrefixOptions) Prefix {
	info.resetCache()
	if options == nil {
		options = &PrefixOptions{}
	}
	options.Defaults()
	result := &prefix{
		info:      info,
		options:   options,
		responses: map[int]*response{},
	}

	// copy responses
	for key, val := range info.responses {
		result.responses[key] = val
	}

	return result
}

// prefix is special path that only creates new paths
type prefix struct {
	info      *prefixInfo
	options   *PrefixOptions
	responses map[int]*response
}

func (p *prefix) Prefix(pathPrefix string, options ...*PrefixOptions) Prefix {
	p.info.resetCache()
	var opts *PrefixOptions
	if len(options) > 0 && options[0] != nil {
		opts = options[0]
	}
	info := *p.info

	// take care about paths
	{
		if !strings.HasSuffix(info.pathPrefix, "/") {
			info.pathPrefix = info.pathPrefix + "/"
		}
		info.pathPrefix = info.pathPrefix + strings.TrimPrefix(pathPrefix, "/")
	}

	return &prefix{
		info:      &info,
		options:   opts,
		responses: map[int]*response{},
	}
}

func (p *prefix) Path(path string, method string, options ...*PathOptions) Path {
	p.info.resetCache()
	var opts *PathOptions
	if len(options) > 0 && options[0] != nil {
		opts = options[0]
	}

	return newPath(&pathInfo{
		Path:        p.info.pathPrefix + path,
		Method:      method,
		Definitions: p.info.definitions,
		Invalidate:  p.info.invalidate,
		Options:     opts,
		Swagger:     p.info.swagger,
	})
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
