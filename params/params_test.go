package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParams(t *testing.T) {
	t.Run("test simple types", func(t *testing.T) {
		t.Run("test invalid types", func(t *testing.T) {
			for _, invalid := range []interface{}{
				1, int8(1), int16(1), int32(1), int64(1),
				uint(1), uint8(1), uint16(1), uint32(1), uint64(1),
				string("oops"), float32(1.0), float64(1.0), complex64(1.0), complex128(1.0),
			} {
				assert.Panics(t, func() {
					New().Add(invalid)
				})
			}
		})

		t.Run("test valid types", func(t *testing.T) {
			type Query struct {
				Value string `json:"value" swag:"description='some value'"`
			}
			p := New()
			p.Add(Query{})
			s := p.Spec()
			_ = s
			assert.Equal(t, 1, len(p.Spec()))
			assert.Equal(t, "value", p.Spec()[0].Name)
			assert.Equal(t, "some value", p.Spec()[0].Description)
			p.Add(nil)
			assert.Equal(t, 0, len(p.Spec()))
		})
	})
}
