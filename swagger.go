package swag

import (
	"encoding/json"
	"net/http"

	"github.com/go-openapi/spec"
	"github.com/matryer/resync"
)

type SwaggerOptions struct {
	Description string
	Version     string
	Host        string
	BasePath    string
	License     *License
	Contact     *ContactInfo
}

// Defaults fill blank values
func (s *SwaggerOptions) Defaults() {

}

// New returns new swagger
func New(title string, options ...*SwaggerOptions) Swagger {
	var opts *SwaggerOptions
	if len(options) > 0 && options[0] != nil {
		opts = options[0]
	} else {
		opts = &SwaggerOptions{}
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
	}
}

type swagger struct {
	spec        spec.Swagger
	options     *SwaggerOptions
	definitions spec.Definitions
	once        resync.Once
	generated   *spec.Swagger
}

// ServeHTTP gives ability to use it in net/http
func (s *swagger) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// add json header
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
		if s.generated, err = s.Spec(); err != nil {
			return
		}
	})

	if err != nil {
		return
	}

	response, err = json.Marshal(s.generated)
	return
}

// Spec returns spec swagger
func (s *swagger) Spec() (*spec.Swagger, error) {
	var paths = spec.Paths{
		VendorExtensible: spec.VendorExtensible{Extensions: map[string]interface{}{"x-framework": "swag-go"}},
		Paths: map[string]spec.PathItem{
			"/": {
				PathItemProps: spec.PathItemProps{
					Get: spec.NewOperation("what").WithTags().WithID("getThing"),
					//Put:        nil,
					//Post:       nil,
					//Delete:     nil,
					//Options:    nil,
					//Head:       nil,
					//Patch:      nil,
					//Parameters: nil,
				},
				//Refable: spec.Refable{Ref: spec.MustCreateRef("cats")},
			},
		},
	}
	s.spec.Paths = &paths

	return &s.spec, nil
}

func (s *swagger) Path(p string, method string, options ...*PathOptions) Path {
	// reset generated thing
	s.once.Reset()

	var opts *PathOptions
	if len(options) > 0 && options[0] != nil {
		opts = options[0]
	}

	return newPath(&pathInfo{
		Path:        p,
		Method:      method,
		Definitions: s.definitions,
		Options:     opts,
		Invalidate:  func() { s.once.Reset() },
	})
}
