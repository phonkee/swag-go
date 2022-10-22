package swag

import (
	"encoding/json"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-openapi/spec"
	"github.com/matryer/resync"
	"github.com/phonkee/swag-go/definitions"
)

// New returns new swag
func New(title string, options ...*Options) Swagger {
	var opts *Options
	if len(options) > 0 && options[0] != nil {
		opts = options[0]
	} else {
		opts = &Options{}
	}
	opts.Defaults()
	return &swag{
		title:       title,
		options:     opts,
		definitions: definitions.New(),
		paths:       make([]*path, 0),
	}
}

// swag implementation of Swagger
type swag struct {
	title         string
	specification spec.Swagger
	options       *Options
	definitions   definitions.Interface
	once          resync.Once
	cached        *spec.Swagger
	paths         []*path
}

func (s *swag) Definitions() definitions.Interface {
	return s.definitions
}

func (s *swag) Debug() {
	for _, p := range s.paths {
		println("path", spew.Sdump(p))
	}
}

func (s *swag) addPath(p *path) {
	s.paths = append(s.paths, p)
}

// MarshalJSON marshals into json and caches result
func (s *swag) MarshalJSON() (response []byte, err error) {
	return json.Marshal(s.spec())
}

// Path returns path
func (s *swag) Path(p string, method string, options ...*PathOptions) Path {
	// reset generated thing
	s.once.Reset()

	var opts *PathOptions
	if len(options) > 0 && options[0] != nil {
		opts = options[0]
	}

	np := newPath(&pathInfo{
		Path:        p,
		Method:      method,
		Definitions: s.definitions,
		Options:     opts,
		Invalidate:  s.invalidate,
		Swagger:     s,
	})

	// add path to swag
	s.addPath(np)

	return np
}

// Prefix returns prefixed prefix
func (s *swag) Prefix(pathPrefix string, options ...*PrefixOptions) Prefix {
	var opts *PrefixOptions
	if len(options) > 0 && options[0] != nil {
		opts = options[0]
	}
	return newPrefix(&prefixInfo{
		definitions: s.definitions,
		swagger:     s,
		pathPrefix:  pathPrefix,
		resetCache: func() {
			s.once.Reset()
		},
		responses:  map[int]*response{},
		invalidate: s.invalidate,
	}, opts)
}

// ServeHTTP gives ability to use it in net/http
func (s *swag) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(s); err != nil {
		http.Error(writer, "cannot encode json", http.StatusInternalServerError)
		return
	}
}

// spec returns specification swag
// TODO: finish this
func (s *swag) spec() *spec.Swagger {
	// only once please
	s.once.Do(func() {
		s.specification = spec.Swagger{
			VendorExtensible: spec.VendorExtensible{},
			SwaggerProps: spec.SwaggerProps{
				//ID:      "http://localhost:3849/api-docs",
				Swagger:  "2.0",
				Consumes: []string{"application/json"},
				Produces: []string{"application/json"},
				Schemes:  []string{"https"},
				Info: &spec.Info{
					InfoProps: spec.InfoProps{
						Description:    s.options.Description,
						Title:          s.title,
						TermsOfService: s.options.TermsOfServices,
						Contact:        s.options.Contact.Spec(),
						License:        s.options.License.Spec(),
						Version:        s.options.Version,
					},
				},
				// Host:     "some.api.out.there",
				// BasePath: "",
				Paths: &spec.Paths{
					VendorExtensible: spec.VendorExtensible{Extensions: map[string]interface{}{"x-framework": XFramework}},
					Paths:            map[string]spec.PathItem{},
				},
			},
		}

		for _, p := range s.paths {

			for k, v := range p.spec().Paths {
				if _, ok := s.specification.Paths.Paths[k]; !ok {
					s.specification.Paths.Paths[k] = spec.PathItem{
						PathItemProps: spec.PathItemProps{
							Parameters: []spec.Parameter{},
						},
					}
				}

				temp := s.specification.Paths.Paths[k]
				if v.Get != nil {
					temp.PathItemProps.Get = v.Get
				}

				s.specification.Paths.Paths[k] = temp
			}
		}
	})
	return &s.specification
}

func (s *swag) invalidate() {
	s.once.Reset()
}
