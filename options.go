package swag

import "strings"

type Options struct {
	Description     string
	Version         string
	Host            string
	License         *License
	Contact         *ContactInfo
	TermsOfServices string
}

// Defaults fill blank values
func (s *Options) Defaults() {
	if s.Version == "" {
		s.Version = DefaultVersion
	}
}

type PathOptions struct {
	Description string
	ID          string
	Tags        []string
	Deprecated  bool
}

func (p *PathOptions) Defaults() {

}

type PrefixOptions struct {
	// always provided description (other will be appended)
	Description string
}

// Defaults sets default values correctly (even fix invalid values)
func (p *PrefixOptions) Defaults() {
}

type ResponseOptions struct {
	Description string
	Deprecated  bool
}

func (r *ResponseOptions) Defaults() {
	if r == nil {
		*r = ResponseOptions{}
	}
	r.Description = strings.TrimSpace(r.Description)
}
