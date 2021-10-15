package swag

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPathParams struct {
	Name string `json:"name" swag_description:"Hello this is description"`
}

type BaseResponse struct {
	Status int `json:"status" swag_description:"This is status from headers"`
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
