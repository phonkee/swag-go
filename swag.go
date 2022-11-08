package swag

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/go-openapi/spec"
	"github.com/matryer/resync"
	"github.com/phonkee/swag-go/definitions"
)

// New returns new swag
func New(title string, options ...*Options) Swag {
	defs := definitions.New()
	return &swag{
		title:       title,
		options:     defaultOptions().Merge(options...),
		definitions: defs,
		responses:   make(Responses),
		updaters:    make([]UpdateSpec, 0),
	}
}

// swag implementation of Swag
type swag struct {
	title         string
	specification spec.Swagger
	options       *Options
	definitions   definitions.Definitions
	once          resync.Once
	cached        *spec.Swagger
	responses     Responses
	updaters      []UpdateSpec
}

// no-op
func (s *swag) Debug() {
}

// MarshalJSON marshals into json and caches result
func (s *swag) MarshalJSON() (response []byte, err error) {
	return json.Marshal(s.spec())
}

// ServeHTTP gives ability to use it in net/http
func (s *swag) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(s); err != nil {
		http.Error(writer, "cannot encode json", http.StatusInternalServerError)
		return
	}
}

// spec returns specification swag
// TODO: finish this
func (s *swag) spec() *spec.Swagger {
	// only once please
	s.once.Do(func() {
		// get specification
		s.specification = defaultSpec(s.title, s.options)

		// updaters other than s.updaters
		updaters := []UpdateSpec{
			s.options.License,
			s.options.Contact,
		}

		for _, updater := range append(updaters, s.updaters...) {
			if err := updater.UpdateSpec(&s.specification); err != nil {
				panic(err)
			}
		}

		if s.definitions.UpdateSpec(&s.specification) != nil {
			panic("cannot update definitions")
		}
	})
	return &s.specification
}

func (s *swag) invalidate() {
	s.once.Reset()
}

// RegisterType registers type with custom marshalling
func (s *swag) RegisterType(what interface{}, fn func(schema *spec.Schema)) {
	s.definitions.RegisterType(reflect.TypeOf(what), fn)
}

func (s *swag) UpdateSpec(swagger *spec.Swagger) error {
	for _, upd := range s.updaters {
		if err := upd.UpdateSpec(swagger); err != nil {
			return err
		}
	}
	return nil
}
