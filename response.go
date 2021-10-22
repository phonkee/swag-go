package swag

type response struct {
	status int
	target interface{}
}

func newResponse(status int, target interface{}) *response {
	return &response{
		status: status,
		target: target,
	}
}
