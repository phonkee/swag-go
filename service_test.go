package swag

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPath struct {
	Name string `json:"name"`
}

type TestResponse struct {
	Status int    `json:"status"`
	Result string `json:"result"`
}

func TestNew(t *testing.T) {
	svc := New("service")
	assert.NotNil(t, svc)

	svc.Path("/hello/world").
		Params(TestPath{}).
		Response(http.StatusOK, TestResponse{})
}
