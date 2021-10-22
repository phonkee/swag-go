package swag

import "github.com/go-openapi/spec"

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
