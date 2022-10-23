package swag

type Responses map[int]*response

func (r Responses) Clone() Responses {
	result := make(Responses, len(r))
	for k, v := range r {
		result[k] = v
	}
	return result
}

type response struct {
	options *ResponseOptions
	status  int
}

func (s *swag) Response(status int, response interface{}, options ...*ResponseOptions) Swag {
	// TODO: add here
	return s
}

func (p *pathImpl) Response(status int, what interface{}, options ...*ResponseOptions) Path {
	//opts := defaultResponseOptions().Merge(options...)
	// TODO: add response here
	return p
}

func (p *prefix) Response(status int, response interface{}, options ...*ResponseOptions) Prefix {
	// TODO: implement me
	return p
}
