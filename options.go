package swag

type Options struct {
	Description     string
	Version         string
	Host            string
	License         *License
	Contact         *ContactInfo
	TermsOfServices string
}

func defaultOptions() *Options {
	return &Options{
		Version: DefaultVersion,
	}
}

func (o *Options) Merge(opts ...*Options) *Options {
	for _, opt := range opts {
		_ = opt
	}
	return o
}

type PathOptions struct {
	Description string
	ID          string
	Tags        []string
	Deprecated  bool
}

func defaultPathOptions() *PathOptions {
	return &PathOptions{}
}

func (p *PathOptions) Merge(others ...*PathOptions) *PathOptions {
	for _, other := range others {
		_ = other
	}
	return p
}

type PrefixOptions struct {
	// always provided description (other will be appended)
	Description string
}

func defaultPrefixOptions() *PrefixOptions {
	return &PrefixOptions{}
}

func (p *PrefixOptions) Clone() *PrefixOptions {
	return &*p
}

func (p *PrefixOptions) Merge(others ...*PrefixOptions) *PrefixOptions {
	for _, other := range others {
		_ = other
	}
	return p
}

type ResponseOptions struct {
	Description string
	Deprecated  bool
}

func defaultResponseOptions() *ResponseOptions {
	return &ResponseOptions{}
}

func (r *ResponseOptions) Merge(others ...*ResponseOptions) *ResponseOptions {
	for _, other := range others {
		_ = other
	}
	return r
}
