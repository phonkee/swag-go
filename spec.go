package swag

import "github.com/go-openapi/spec"

func defaultSpec(title string, opts *Options) spec.Swagger {
	o := defaultOptions().Merge(opts)

	return spec.Swagger{
		VendorExtensible: spec.VendorExtensible{},
		SwaggerProps: spec.SwaggerProps{
			//ID:      "http://localhost:3849/api-docs",
			Swagger:  "2.0",
			Consumes: []string{"application/json"},
			Produces: []string{"application/json"},
			Schemes:  []string{"https"},
			Info: &spec.Info{
				InfoProps: spec.InfoProps{
					Description:    o.Description,
					Title:          title,
					TermsOfService: o.TermsOfServices,
					Version:        o.Version,
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
}
