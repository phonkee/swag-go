package swag

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPathParams struct {
	Name string `json:"name"`
}

type BaseResponse struct {
	Status int `json:"status"`
}

type TestResponse struct {
	BaseResponse
}

func TestNew(t *testing.T) {
	svc := New("service")
	assert.NotNil(t, svc)

	// add post method
	svc.Path("/hello/world", http.MethodPost).
		Params(TestPathParams{}).
		Response(http.StatusOK, TestResponse{}).
		Response(http.StatusNotFound, BaseResponse{}).
		Response(http.StatusInternalServerError, BaseResponse{})
}
