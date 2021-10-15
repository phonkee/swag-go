package swag

import (
	"github.com/matryer/resync"
)

// New returns new service
func New(name string, options ...*ServiceOptions) Service {
	return &service{}
}

type service struct {
	once resync.Once
}

// MarshalJSON marshals into json and caches result
func (s *service) MarshalJSON() ([]byte, error) {
	panic("implement me")
}

func (s *service) Path(path string, method string, options ...*PathOptions) Path {
	s.once.Reset()
	return nil
}
