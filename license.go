package swag

import "github.com/go-openapi/spec"

type License struct {
	Name string
	URL  string
}

func (l *License) UpdateSpec(swagger *spec.Swagger) error {
	if l == nil {
		return nil
	}
	if swagger.SwaggerProps.Info == nil {
		return nil
	}
	swagger.SwaggerProps.Info.License = &spec.License{
		LicenseProps: spec.LicenseProps{
			Name: l.Name,
			URL:  l.URL,
		},
	}
	return nil
}
