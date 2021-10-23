package swag

import (
	"encoding/json"
	"net/http"

	"github.com/go-openapi/spec"
	"github.com/matryer/resync"
)

type Options struct {
	Description string
	Version     string
	Host        string
	BasePath    string
	License     *License
	Contact     *ContactInfo
}

// Defaults fill blank values
func (s *Options) Defaults() {

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
		spec: spec.Swagger{
			VendorExtensible: spec.VendorExtensible{},
			SwaggerProps: spec.SwaggerProps{
				//ID:      "http://localhost:3849/api-docs",
				Swagger:  "2.0",
				Consumes: []string{"application/json"},
				Produces: []string{"application/json"},
				Schemes:  []string{"http", "https"},
				Info: &spec.Info{
					InfoProps: spec.InfoProps{
						Description: opts.Description,
						Title:       title,
						//TermsOfService: "",
						Contact: opts.Contact.Spec(),
						License: opts.License.Spec(),
						Version: opts.Version,
					},
				},
				Host:     "some.api.out.there",
				BasePath: "/",
			},
		},
		options:     opts,
		definitions: make(spec.Definitions),
		paths:       make([]*path, 0),
	}
}

type swagger struct {
	spec        spec.Swagger
	options     *Options
	definitions spec.Definitions
	once        resync.Once
	generated   *spec.Swagger
	paths       []*path
}

func (s *swagger) Prefix(pathPrefix string) Prefix {
	return newPrefix(&prefixInfo{
		swagger:    s,
		pathPrefix: pathPrefix,
	})
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

// MarshalJSON marshals into json and caches result
func (s *swagger) MarshalJSON() (response []byte, err error) {
	s.once.Do(func() {
		// generate here
		s.generated = s.Spec()
	})

	if err != nil {
		return
	}

	response, err = json.Marshal(s.generated)
	return
}

// Spec returns spec swagger
func (s *swagger) Spec() *spec.Swagger {
	var paths = spec.Paths{
		VendorExtensible: spec.VendorExtensible{Extensions: map[string]interface{}{"x-framework": XFramework}},
		Paths:            map[string]spec.PathItem{
			//"/": {
			//	PathItemProps: spec.PathItemProps{
			//		Get: spec.NewOperation("what").WithTags().WithID("getThing"),
			//		//Put:        nil,
			//		//Post:       nil,
			//		//Delete:     nil,
			//		//Options:    nil,
			//		//Head:       nil,
			//		//Patch:      nil,
			//		//Parameters: nil,
			//	},
			//	//Refable: spec.Refable{Ref: spec.MustCreateRef("cats")},
			//},
		},
	}
	s.spec.Paths = &paths

	for _, p := range s.paths {
		for k, v := range p.Spec().Paths {
			if _, ok := s.spec.Paths.Paths[k]; !ok {
				s.spec.Paths.Paths[k] = spec.PathItem{
					PathItemProps: spec.PathItemProps{
						Parameters: []spec.Parameter{},
					},
				}
			}

			where := s.spec.Paths.Paths[k]

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

	return &s.spec
}

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
		Invalidate:  func() { s.once.Reset() },
		Swagger:     s,
	})
	s.paths = append(s.paths, np)

	return np
}
