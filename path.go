package swag

import (
	"github.com/go-openapi/spec"
	"github.com/phonkee/swag-go/definitions"
	"github.com/phonkee/swag-go/params"
)

func (s *swag) Path(path string, method string, options ...*PathOptions) Path {
	result := &pathImpl{
		path:        path,
		definitions: s.definitions,
		responses:   s.responses.Clone(),
		pathParams:  params.New(),
		queryParams: params.New(),
		method:      method,
		options:     defaultPathOptions().Merge(options...),
		invalidate:  s.invalidate,
	}
	s.updaters = append(s.updaters, result)

	return result
}

func (p *prefix) Path(path string, method string, options ...*PathOptions) Path {
	result := &pathImpl{
		path:        pathJoin(p.path, path),
		definitions: p.definitions,
		responses:   p.responses.Clone(),
		pathParams:  p.pathParams.Clone(),
		queryParams: p.queryParams.Clone(),
		method:      method,
		options:     defaultPathOptions().Merge(options...),
		invalidate:  p.invalidate,
	}
	p.updaters = append(p.updaters, result)

	return result
}

// pathImpl implementation
type pathImpl struct {
	path        string
	definitions definitions.Definitions
	responses   Responses
	pathParams  params.Params
	queryParams params.Params
	method      string
	options     *PathOptions
	invalidate  func()
}

func (p *pathImpl) Body(i interface{}) Path {
	defer p.invalidate()
	//TODO implement me
	panic("implement me")
}

func (p *pathImpl) PathParams(i interface{}) Path {
	defer p.invalidate()
	p.pathParams.Add(i)
	return p
}

func (p *pathImpl) QueryParams(i interface{}) Path {
	defer p.invalidate()
	p.queryParams.Add(i)
	return p
}

func (p *pathImpl) UpdateSpec(swagger *spec.Swagger) error {
	// TODO: go over all responses and add them to swagger

	//for _, p := range s.paths {
	//	_ = p
	//
	//	//for k, v := range p.spec().Paths {
	//	//	if _, ok := s.specification.Paths.Paths[k]; !ok {
	//	//		s.specification.Paths.Paths[k] = spec.PathItem{
	//	//			PathItemProps: spec.PathItemProps{
	//	//				Parameters: []spec.Parameter{},
	//	//			},
	//	//		}
	//	//	}
	//	//
	//	//	temp := s.specification.Paths.Paths[k]
	//	//	if v.Get != nil {
	//	//		temp.PathItemProps.Get = v.Get
	//	//	}
	//	//
	//	//	s.specification.Paths.Paths[k] = temp
	//	//}
	//}

	return nil
}
