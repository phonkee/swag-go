package params

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/fatih/structs"
	"github.com/go-openapi/spec"
	"github.com/phonkee/swag-go/utils"
)

type Params interface {
	// Add adds a parameter to the list of parameters, if nil is provided it removes all previous parameter
	Add(value interface{})
	Spec() []spec.Parameter
}

func New() Params {
	return &params{
		params: make([]spec.Parameter, 0),
	}
}

type params struct {
	params []spec.Parameter
}

func (p *params) Add(value interface{}) {
	if value == nil {
		p.params = make([]spec.Parameter, 0)
		return
	}
	typ := reflect.TypeOf(value)
	for {
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
			continue
		} else {
			break
		}
	}

	if typ.Kind() != reflect.Struct {
		panic(fmt.Sprintf("Params.Add accepts only struct, got %s", typ.Kind()))
	}

	for _, field := range structs.New(value).Fields() {
		pt, err := paramTag.ParseTag(field.Tag("swag"))
		if err != nil {
			panic(fmt.Sprintf("error parsing tag %s", err))
		}

		// prepare param
		param := spec.Parameter{}
		param.Description = pt.Description
		param.Name = utils.GetFieldName(field)

		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			// mixing int and uint
			param.Type = "integer"
			switch typ.Kind() {
			case reflect.Int, reflect.Uint:
				param.Format = "int" + strconv.FormatInt(int64(typ.Bits()), 10)
			case reflect.Int8, reflect.Uint8:
				param.Format = "int8"
			case reflect.Int16, reflect.Uint16:
				param.Format = "int16"
			case reflect.Int32, reflect.Uint32:
				param.Format = "int32"
			case reflect.Int64, reflect.Uint64:
				param.Format = "int64"
			}
		case reflect.String:
			param.Type = "string"
		default:
			panic(fmt.Sprintf("don't know how to handle this type: %v", field))
		}
		p.params = append(p.params, param)
	}

}

// Spec returns swagger params
func (p *params) Spec() []spec.Parameter {
	return p.params
}
