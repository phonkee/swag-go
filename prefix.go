package swag

import (
	"strings"

	"github.com/go-openapi/spec"
)

// TODO: prefix is not implemented yet, should not be very complicated though
// Currently it's lower priority but soon it will be implemented

type PrefixOptions struct {
	// always provided description (other will be appended)
	Description string
}

// Defaults sets default values correctly (even fix invalid values)
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
	// fix defaults
	options.Defaults()

	// prepare result
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
	return p.copy(pathPrefix, options...)
}

func (p *prefix) Path(path string, method string, options ...*PathOptions) Path {
	p.info.resetCache()
	var opts *PathOptions
	if len(options) > 0 && options[0] != nil {
		opts = options[0]
	} else {
		opts = &PathOptions{}
	}
	opts.Defaults()

	result := newPath(&pathInfo{
		Path:        p.info.pathPrefix + path,
		Method:      method,
		Definitions: p.info.definitions,
		Invalidate:  p.info.invalidate,
		Options:     opts,
		Swagger:     p.info.swagger,
	})

	// add path
	p.info.swagger.addPath(result)

	return result
}

func (p *prefix) Response(status int, response interface{}, options ...*ResponseOptions) Prefix {
	p.info.resetCache()
	var opts *ResponseOptions
	if len(options) > 0 && options[0] != nil {
		opts = options[0]
	} else {
		opts = &ResponseOptions{}
	}
	// set defaults
	opts.Defaults()

	// overwrite old response
	p.responses[status] = newResponse(status, response, opts)
	return p
}

func (p *prefix) PathParams(i interface{}) Prefix {
	p.info.resetCache()
	//for _, param := range p.getParams(i, ParamTypePath) {
	//	p.item.PathItemProps.Parameters = append(p.item.PathItemProps.Parameters, *param)
	//}
	return p
}

func (p *prefix) QueryParams(i interface{}) Prefix {
	p.info.resetCache()
	//for _, param := range p.getParams(i, ParamTypePath) {
	//	p.item.PathItemProps.Parameters = append(p.item.PathItemProps.Parameters, *param)
	//}
	return p
}

func (p *prefix) copy(pathPrefix string, options ...*PrefixOptions) Prefix {
	p.info.resetCache()
	opts := &*p.options
	info := *p.info

	// take care about paths
	{
		if !strings.HasSuffix(info.pathPrefix, "/") {
			info.pathPrefix = info.pathPrefix + "/"
		}
		info.pathPrefix = info.pathPrefix + strings.TrimPrefix(pathPrefix, "/")
	}

	result := &prefix{
		info:      &info,
		options:   opts,
		responses: map[int]*response{},
	}

	// copy all responses
	for k, v := range p.responses {
		result.responses[k] = v
	}

	return result

}
