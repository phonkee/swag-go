package swag

import (
	"github.com/matryer/resync"
)

// New returns new service
func New(name string, options ...*ServiceOptions) Service {
	return &service{
		components: newComponents(),
		generated:  []byte(""),
	}
}

type service struct {
	components Components
	once       resync.Once
	generated  []byte
}

// MarshalJSON marshals into json and caches result
func (s *service) MarshalJSON() ([]byte, error) {
	s.once.Do(func() {
		// generate here
	})

	return s.generated, nil
}

func (s *service) Path(p string, method string, options ...*PathOptions) Path {
	var opts *PathOptions
	if len(options) > 0 && options[0] != nil {
		opts = options[0]
	}

	// reset generated
	s.once.Reset()

	return newPath(p, method, s.components, opts)
}
