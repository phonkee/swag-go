package swag

import (
	"github.com/go-openapi/spec"
	"github.com/phonkee/swag-go/definitions"
	"github.com/phonkee/swag-go/params"
)

type prefix struct {
	pathParams  params.Params
	queryParams params.Params
	path        string
	options     *PrefixOptions
	definitions definitions.Definitions
	responses   Responses
	updaters    []UpdateSpec
}

func (s *swag) Prefix(path string, options ...*PrefixOptions) Prefix {
	result := &prefix{
		pathParams:  params.New(),
		queryParams: params.New(),
		path:        path,
		options:     defaultPrefixOptions().Merge(options...),
		definitions: s.definitions,
		responses:   s.responses.Clone(),
		updaters:    make([]UpdateSpec, 0),
	}

	s.updaters = append(s.updaters, result)

	return result
}

func (p *prefix) Prefix(path string, options ...*PrefixOptions) Prefix {
	// TODO: add proper path joining here
	result := &prefix{
		path:        p.path + path,
		pathParams:  p.pathParams.Clone(),
		queryParams: p.queryParams.Clone(),
		options:     p.options.Clone().Merge(options...),
		definitions: p.definitions,
		responses:   p.responses.Clone(),
		updaters:    make([]UpdateSpec, 0),
	}

	p.updaters = append(p.updaters, result)

	return result
}

func (p *prefix) PathParams(i interface{}) Prefix {
	p.pathParams.Add(i)
	return p
}

func (p *prefix) QueryParams(i interface{}) Prefix {
	p.queryParams.Add(i)
	return p
}

func (p *prefix) UpdateSpec(swagger *spec.Swagger) error {
	for _, upd := range p.updaters {
		if err := upd.UpdateSpec(swagger); err != nil {
			return err
		}
	}
	return nil
}
