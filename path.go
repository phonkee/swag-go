package swag

import (
	"fmt"
	"net/http"
	"reflect"
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
	Swagger     *swagger
	responses   map[int]*response
}

func newPath(info *pathInfo) *path {
	result := &path{
		info:      info,
		responses: map[int]*response{},
		item: spec.PathItem{
			PathItemProps: spec.PathItemProps{
				Parameters: []spec.Parameter{},
			},
		},
	}

	for k, v := range info.responses {
		result.responses[k] = v
	}

	return result
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
	for _, param := range p.getParams(i, spec.PathParam) {
		p.item.PathItemProps.Parameters = append(p.item.PathItemProps.Parameters, *param)
	}

	return p
}

// QueryParams adds query params
func (p *path) QueryParams(i interface{}) Path {
	p.info.Invalidate()
	for _, param := range p.getParams(i, spec.PathParam) {
		_ = param
	}
	return p
}

// Response adds response to path
func (p *path) Response(status int, what interface{}, options ...*ResponseOptions) Path {
	p.info.Invalidate()

	var opts *ResponseOptions

	if len(options) > 0 && options[0] != nil {
		opts = options[0]
	} else {
		opts = &ResponseOptions{}
	}
	opts.Defaults()

	// no response
	if what == nil {
		p.responses[status] = nil
		return p
	}

	// TODO: when what is nil, we should empty responses?
	p.responses[status] = newResponse(status, what, opts)
	return p
}

// getParams returns params for given struct, typ
func (p *path) getParams(i interface{}, initFunc func(string) *spec.Parameter) []*spec.Parameter {
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

		format := initFunc(name).WithDescription(description)

		// TODO: here comes parameters.go implementation
		var (
			err error
			nf  *spec.Parameter
		)
		if nf, err = getParameter(reflect.TypeOf(i), format); err != nil {
			if err == errParameterNotFound {
				err = nil
			} else {
				panic(err)
			}
		} else {
			format = nf
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

func (p *path) spec() spec.Paths {

	// now add all responses to item
	// TODO: now we need to merge here things, and return correct things

	result := spec.Paths{
		Paths: map[string]spec.PathItem{
			p.info.Path: p.item,
		},
	}

	// TODO: finish this huge thing
	switch p.info.Method {
	case http.MethodGet:
		where := result.Paths[p.info.Path]
		where.PathItemProps.Get = p.operation()
		result.Paths[p.info.Path] = where
	}

	return result
}

func (p *path) operation() *spec.Operation {
	result := &spec.Operation{
		OperationProps: spec.OperationProps{
			Description: p.info.Options.Description,
			//Consumes:     nil,
			//Produces:     nil,
			//Schemes:      nil,
			Tags: p.info.Options.Tags,
			//Summary:      "",
			//ExternalDocs: nil,
			//ID:           "",
			//Deprecated:   false,
			//Security:     nil,
			//Parameters:   nil,
			Responses: &spec.Responses{
				ResponsesProps: spec.ResponsesProps{
					StatusCodeResponses: map[int]spec.Response{},
				},
			},
		},
	}
	for status, response := range p.responses {
		_ = response
		result.OperationProps.Responses.ResponsesProps.StatusCodeResponses[status] = spec.Response{}
	}

	// TODO: add all parameters and responses

	return result
}
