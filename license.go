package swag

import "github.com/go-openapi/spec"

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
