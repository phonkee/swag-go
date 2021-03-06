package swag

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-openapi/spec"
)

type paramType string

const (
	ParamTypeQuery paramType = "query"
	ParamTypePath  paramType = "path"
)

type blankResponse int

type pathInfo struct {
	Path        string
	Method      string
	Definitions spec.Definitions
	Options     *PathOptions
	Invalidate  func()
	Swagger     *swagger
}

func newPath(info *pathInfo) *path {
	return &path{
		info:      info,
		responses: map[int]*response{},
		item: spec.PathItem{
			PathItemProps: spec.PathItemProps{
				Parameters: []spec.Parameter{},
			},
		},
	}
}

type path struct {
	responses map[int]*response
	item      spec.PathItem
	info      *pathInfo
}

func (p *path) Body(i interface{}) Path {
	p.info.Invalidate()
	return p
}

// PathParams adds path params
func (p *path) PathParams(i interface{}) Path {
	p.info.Invalidate()
	for _, param := range p.Params(i, ParamTypePath) {
		p.item.PathItemProps.Parameters = append(p.item.PathItemProps.Parameters, *param)
	}

	return p
}

// QueryParams adds query params
func (p *path) QueryParams(i interface{}) Path {
	p.info.Invalidate()
	for _, param := range p.Params(i, ParamTypeQuery) {
		_ = param
		// spew.Dump(param)
	}
	return p
}

func (p *path) Params(i interface{}, typ paramType) []*spec.Parameter {
	result := make([]*spec.Parameter, 0)
	ss := structs.New(i)
	for index, field := range ss.Fields() {
		description := getFieldDescription(field)
		name := field.Name()
		if jsonName := field.Tag("json"); jsonName != "" {
			splitted := strings.SplitN(jsonName, ",", 2)
			if jsonName = strings.TrimSpace(splitted[0]); jsonName != "" {
				name = jsonName
			}
		}

		var format *spec.Parameter

		switch typ {
		case ParamTypeQuery:
			format = spec.QueryParam(name).WithDescription(description)
		case ParamTypePath:
			format = spec.PathParam(name).WithDescription(description)
		}

		// get kind
		kind := field.Kind()

		if kind == reflect.Ptr {
			format.Required = false
			// not nice, but not accessible from field
			kind = reflect.TypeOf(i).Field(index).Type.Elem().Kind()
		}

		var typ, tmpFmt string
		// now type switch for types
		// TODO: finish https://github.com/OAI/OpenAPI-Specification/blob/main/versions/2.0.md#data-types
		switch kind {
		// simplified (no format)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			typ = "integer"
		case reflect.Float32:
			typ = "number"
			tmpFmt = "float"
		case reflect.Float64:
			typ = "number"
			tmpFmt = "double"
		case reflect.String:
			typ = "string"
		case reflect.Bool:
			typ = "boolean"
		default:
			panic(fmt.Sprintf("unsupported kind %v", kind.String()))
		}

		format.SimpleSchema.Type = typ
		format.SimpleSchema.Format = tmpFmt

		p.item.PathItemProps.Parameters = append(p.item.PathItemProps.Parameters, *format)

		result = append(result, format)
	}

	return result
}

// Response adds response to path
func (p *path) Response(status int, what interface{}, options ...*ResponseOptions) Path {
	p.info.Invalidate()

	var opts *ResponseOptions

	if len(options) > 0 {
		opts = options[0]
	}

	// no response
	if what == nil {
		p.responses[status] = nil
		return p
	}

	// TODO: when what is nil, we should empty responses?
	p.responses[status] = newResponse(status, what, opts)
	return p
}

func (p *path) Spec() spec.Paths {

	// now add all responses to item
	// TODO: now we need to merge here things, and return correct things

	result := spec.Paths{
		Paths: map[string]spec.PathItem{
			p.info.Path: p.item,
		},
	}

	return result
}
