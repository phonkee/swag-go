package swag

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-openapi/spec"
)

var (
	// ParamsScopeDelimiter is string how we join parameters such as (user.id), you can change it if you want
	ParamsScopeDelimiter = "."
)

// inspectParams inspects given target and calls fn callback on each parameter
func inspectParams(target interface{}, fn func(name string) *spec.Parameter) []*spec.Parameter {
	result := make([]*spec.Parameter, 0)
	ss := structs.New(target)
	for index, field := range ss.Fields() {
		description := getFieldDescription(field)
		name := field.Name()
		if jsonName := field.Tag("json"); jsonName != "" {
			splitted := strings.SplitN(jsonName, ",", 2)
			if jsonName = strings.TrimSpace(splitted[0]); jsonName != "" {
				name = jsonName
			}
		}

		format := fn(name).WithDescription(description)

		var (
			err error
			nf  *spec.Parameter
		)

		if nf, err = globalParameters.Get(field.Value(), format); err == nil {
			result = append(result, nf)
			continue
		} else {
			if !errors.Is(err, errParameterNotFound) {
				panic(err)
			}
		}

		// get kind
		kind := field.Kind()

		if kind == reflect.Ptr {
			format.Required = false
			// not nice, but not accessible from field
			kind = reflect.TypeOf(target).Field(index).Type.Elem().Kind()
		}

		// now type switch for kinds
		// types are defined in parameters.go
		// TODO: finish https://github.com/OAI/OpenAPI-Specification/blob/main/versions/2.0.md#data-types
		switch kind {
		// simplified (no format)
		case reflect.Struct:
			for _, subParam := range inspectParams(field.Value(), fn) {
				// only when field is not embedded we add scope
				if !field.IsEmbedded() {
					subParam.Name = fmt.Sprintf("%v%v%v", name, ParamsScopeDelimiter, subParam.Name)
				}
				result = append(result, subParam)
			}
			continue
		case reflect.Slice, reflect.Array:
			// TODO: finish slice/array support
			format.Type = "array"
			format.Items = &spec.Items{
				SimpleSchema: spec.SimpleSchema{
					Type: "__notimplementednow__",
				},
			}

			result = append(result, format)
		default:
			panic(fmt.Sprintf("unsupported kind %v", kind.String()))
		}
	}

	return result
}

// inspectSchema inspects target and returns Schema
func inspectSchema(target interface{}, defs spec.Definitions) (result *spec.Schema) {
	var err error

	if result, err = getSchema(target, defs); err != nil {
		panic(fmt.Sprintf("cannot inspect %T: %v", target, err))
	}

	return
}
