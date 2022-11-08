package swag

import "github.com/go-openapi/spec"

type ContactInfo struct {
	Name  string
	URL   string
	Email string
}

func (c *ContactInfo) UpdateSpec(swagger *spec.Swagger) error {
	if c == nil {
		return nil
	}
	if swagger.SwaggerProps.Info == nil {
		return nil
	}
	swagger.SwaggerProps.Info.Contact = &spec.ContactInfo{
		ContactInfoProps: spec.ContactInfoProps{
			Name:  c.Name,
			URL:   c.URL,
			Email: c.Email,
		},
	}
	return nil
}
