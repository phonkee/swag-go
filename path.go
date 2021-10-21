package swag

import (
	"strings"

	"github.com/fatih/structs"
	"github.com/go-openapi/spec"
)

type pathInfo struct {
	Path        string
	Method      string
	Definitions spec.Definitions
	Options     *PathOptions
	Invalidate  func()
}

func newPath(info *pathInfo) *path {
	return &path{
		definitions: info.Definitions,
		path:        info.Path,
		method:      info.Method,
		options:     info.Options,
		responses:   map[int]*response{},
		invalidate:  info.Invalidate,
		item: spec.PathItem{
			PathItemProps: spec.PathItemProps{
				Parameters: []spec.Parameter{},
			},
		},
	}
}

type path struct {
	path        string
	method      string
	definitions spec.Definitions
	options     *PathOptions
	responses   map[int]*response
	invalidate  func()
	item        spec.PathItem
}

func (p *path) Body(i interface{}) Path {
	p.invalidate()
	return p
}

type response struct {
	schema *spec.Schema
}

func (p *path) PathParams(i interface{}) Path {
	p.invalidate()
	ss := structs.New(i)
	for _, field := range ss.Fields() {
		description := field.Tag("swag_description")
		name := field.Name()
		if jsonName := field.Tag("json"); jsonName != "" {
			splitted := strings.SplitN(jsonName, ",", 2)
			if jsonName = strings.TrimSpace(splitted[0]); jsonName != "" {
				name = jsonName
			}
		}

		_, _ = name, description
		//println(name, field.Kind().String())
		//
		//println(description)
	}

	//println(ss.Name())
	//_ = ss

	//ref := spec.MustCreateRef(fmt.Sprintf("#/definitions/%v", ss.Name()))
	//_ = ref
	return p
}

func (p *path) QueryParams(i interface{}) Path {
	p.invalidate()
	return p
}

func (p *path) Response(status int, what interface{}) Path {
	p.invalidate()
	//_ = p.components.GetSchema(what)

	return p
}

func (p *path) Spec() (*spec.Paths, error) {
	return nil, nil
}
