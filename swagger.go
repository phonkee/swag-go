package swag

import (
	"encoding/json"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-openapi/spec"
	"github.com/matryer/resync"
)

type Options struct {
	Description string
	Version     string
	Host        string
	License     *License
	Contact     *ContactInfo
}

// Defaults fill blank values
func (s *Options) Defaults() {
	if s.Version == "" {
		s.Version = DefaultVersion
	}
}

// New returns new swagger
func New(title string, options ...*Options) Swagger {
	var opts *Options
	if len(options) > 0 && options[0] != nil {
		opts = options[0]
	} else {
		opts = &Options{}
	}
	opts.Defaults()
	return &swagger{
		title:       title,
		options:     opts,
		definitions: make(spec.Definitions),
		paths:       make([]*path, 0),
	}
}

// swagger implementation of Swagger
type swagger struct {
	title         string
	specification spec.Swagger
	options       *Options
	definitions   spec.Definitions
	once          resync.Once
	cached        *spec.Swagger
	paths         []*path
}

func (s *swagger) Debug() {
	for _, p := range s.paths {
		println("path", spew.Sdump(p))
	}
}

func (s *swagger) addPath(p *path) {
	s.paths = append(s.paths, p)
}

// MarshalJSON marshals into json and caches result
func (s *swagger) MarshalJSON() (response []byte, err error) {
	return json.Marshal(s.spec())
}

// Path returns path
func (s *swagger) Path(p string, method string, options ...*PathOptions) Path {
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

	// add path to swagger
	s.addPath(np)

	return np
}

// Prefix returns prefixed prefix
func (s *swagger) Prefix(pathPrefix string, options ...*PrefixOptions) Prefix {
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
func (s *swagger) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(s); err != nil {
		http.Error(writer, "cannot encode json", http.StatusInternalServerError)
		return
	}
}

// spec returns specification swagger
// TODO: finish this
func (s *swagger) spec() *spec.Swagger {
	// only once please
	s.once.Do(func() {
		s.specification = spec.Swagger{
			VendorExtensible: spec.VendorExtensible{},
			SwaggerProps: spec.SwaggerProps{
				//ID:      "http://localhost:3849/api-docs",
				Swagger:  "2.0",
				Consumes: []string{"application/json"},
				Produces: []string{"application/json"},
				Schemes:  []string{"http", "https"},
				Info: &spec.Info{
					InfoProps: spec.InfoProps{
						Description: s.options.Description,
						Title:       s.title,
						//TermsOfService: "",
						Contact: s.options.Contact.Spec(),
						License: s.options.License.Spec(),
						Version: s.options.Version,
					},
				},
				Host:     "some.api.out.there",
				BasePath: "/",
				Paths: &spec.Paths{
					VendorExtensible: spec.VendorExtensible{Extensions: map[string]interface{}{"x-framework": XFramework}},
					Paths:            map[string]spec.PathItem{},
				},
			},
		}

		for _, p := range s.paths {
			spew.Dump(p)
			for k, v := range p.spec().Paths {
				if _, ok := s.specification.Paths.Paths[k]; !ok {
					s.specification.Paths.Paths[k] = spec.PathItem{
						PathItemProps: spec.PathItemProps{
							Parameters: []spec.Parameter{},
						},
					}
				}

				where := s.specification.Paths.Paths[k]

				for _, def := range []struct {
					from *spec.Operation
					to   *spec.Operation
				}{
					{v.Get, where.Get},
					{v.Post, where.Post},
					{v.Put, where.Put},
					{v.Patch, where.Patch},
					{v.Options, where.Options},
					{v.Delete, where.Delete},
					{v.Head, where.Head},
				} {
					if def.from != nil {
						def.to = def.from
					}
				}
			}
		}
	})
	return &s.specification
}

func (s *swagger) invalidate() {
	s.once.Reset()
}
