package swag

import (
	fmt2 "fmt"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/structs"
	"github.com/go-openapi/spec"
)

type pathInfo struct {
	Path        string
	Method      string
	Definitions spec.Definitions
	Options     *PathOptions
	Invalidate  func()
	Swagger     *spec.Swagger
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
		swagger: info.Swagger,
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
	swagger     *spec.Swagger
}

func (p *path) Body(i interface{}) Path {
	p.invalidate()
	return p
}

type response struct {
	schema *spec.Schema
}

// PathParams adds path params
func (p *path) PathParams(i interface{}) Path {
	p.invalidate()
	for _, param := range p.Params(i, spec.PathParam) {
		p.item.PathItemProps.Parameters = append(p.item.PathItemProps.Parameters, *param)
	}

	return p
}

// QueryParams adds query params
func (p *path) QueryParams(i interface{}) Path {
	p.invalidate()
	for _, param := range p.Params(i, spec.QueryParam) {
		spew.Dump(param)
	}
	return p
}

func (p *path) Params(i interface{}, nf func(name string) *spec.Parameter) []*spec.Parameter {
	result := make([]*spec.Parameter, 0)
	ss := structs.New(i)
	for index, field := range ss.Fields() {
		description := field.Tag("swag_description")
		name := field.Name()
		if jsonName := field.Tag("json"); jsonName != "" {
			splitted := strings.SplitN(jsonName, ",", 2)
			if jsonName = strings.TrimSpace(splitted[0]); jsonName != "" {
				name = jsonName
			}
		}
		param := nf(name).WithDescription(description)

		// get kind
		kind := field.Kind()

		if kind == reflect.Ptr {
			param.Required = false
			// not nice, but not accessible from field
			kind = reflect.TypeOf(i).Field(index).Type.Elem().Kind()
		}

		var typ, fmt string
		// now type switch for types
		// TODO: finish https://github.com/OAI/OpenAPI-Specification/blob/main/versions/2.0.md#data-types
		switch kind {
		// simplified (no format)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			typ = "integer"
		case reflect.Float32:
			typ = "number"
			fmt = "float"
		case reflect.Float64:
			typ = "number"
			fmt = "double"
		case reflect.String:
			typ = "string"
		default:
			panic(fmt2.Sprintf("unsupported kind %v", kind.String()))
		}

		param.SimpleSchema.Type = typ
		param.SimpleSchema.Format = fmt

		result = append(result, param)
	}

	return result
}

func (p *path) Response(status int, what interface{}) Path {
	p.invalidate()
	//_ = p.components.GetSchema(what)
	return p
}

func (p *path) Spec() spec.Paths {
	return spec.Paths{
		Paths: map[string]spec.PathItem{},
	}
}
