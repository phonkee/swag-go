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

		if nf, err = globalParameters.Get(reflect.TypeOf(field.Value()), format); err == nil {
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

		// now type switch for types
		// TODO: finish https://github.com/OAI/OpenAPI-Specification/blob/main/versions/2.0.md#data-types
		// TODO: add support for arrays? Is it needed?
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
		default:
			panic(fmt.Sprintf("unsupported kind %v", kind.String()))
		}
	}

	return result
}

// inspectSchema inspects target and returns Schema
func inspectSchema(target interface{}, defs spec.Definitions) (result *spec.Schema) {
	required := true
	typ := reflect.TypeOf(target)

	if typ.Kind() == reflect.Ptr {
		required = false
		typ = typ.Elem()
	}

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:     []string{"integer"},
				Nullable: !required,
			},
		}
	case reflect.String:
		return &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:     []string{"string"},
				Nullable: !required,
			},
		}
	case reflect.Bool:
		return &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:     []string{"boolean"},
				Nullable: !required,
			},
		}
	case reflect.Struct:
		id := fmt.Sprintf("%T", target)
		id = strings.TrimLeft(id, "*")
		if r, ok := defs[id]; ok {
			return &r
		}

		result = spec.RefSchema(id)
		result.ID = id
		result.Properties = spec.SchemaProperties{}
		result.Nullable = !required
		for _, field := range structs.New(target).Fields() {
			if !isFieldAvailable(field) {
				continue
			}
			fsch := inspectSchema(field.Value(), defs)
			name := getFieldName(field)
			if fsch != nil {
				fsch.ID = name
				fsch.Description = getFieldDescription(field)
				result.Properties[name] = *fsch
			}
		}
		defs[id] = *result
		return result
	case reflect.Slice, reflect.Array:
		sch := &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"array"},
			},
		}

		elem := reflect.TypeOf(target).Elem()
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		sch.Items = &spec.SchemaOrArray{
			Schema: inspectSchema(reflect.New(elem).Interface(), defs),
		}

		return sch
	}

	return
}
