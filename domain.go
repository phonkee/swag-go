package swag

import (
	"encoding/json"
	"net/http"

	"github.com/go-openapi/spec"
)

type ContactInfo struct {
	Name  string
	URL   string
	Email string
}

func (c *ContactInfo) Spec() *spec.ContactInfo {
	if c == nil {
		return nil
	}
	return &spec.ContactInfo{
		ContactInfoProps: spec.ContactInfoProps{
			Name:  c.Name,
			URL:   c.URL,
			Email: c.Email,
		},
	}
}

type License struct {
	Name string
	URL  string
}

func (l *License) Spec() *spec.License {
	if l == nil {
		return nil
	}
	return &spec.License{
		LicenseProps: spec.LicenseProps{
			Name: l.Name,
			URL:  l.URL,
		},
	}
}

type Swagger interface {
	http.Handler
	json.Marshaler

	// Path adds new endpoints
	Path(path string, method string, options ...*PathOptions) Path
}

type PathOptions struct {
	Description string
	ID          string
}

type Path interface {
	// Body is request body
	Body(interface{}) Path

	// PathParams adds path params
	PathParams(interface{}) Path

	// QueryParams params
	QueryParams(interface{}) Path

	// Response returned for given status code
	Response(status int, what interface{}) Path

	Spec() spec.Paths
}
