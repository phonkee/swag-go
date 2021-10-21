package swag

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPathParams struct {
	Name     string  `json:"name" swag_description:"Hello this is description"`
	Optional *string `json:"optional"`
}

type TestQueryParams struct {
	Name string `json:"name" swag_description:"Hello this is description"`
}

type BaseResponse struct {
	Status int `json:"status" swag_description:"This is status from headers"`
}

type TestResponse struct {
	BaseResponse
}

func TestNew(t *testing.T) {
	swg := New("Pet store", &Options{
		Description: "Pet store swagger implementation",
		Version:     "1.0.0",
		License: &License{
			Name: "MIT",
		},
	})
	assert.NotNil(t, swg)

	// add post method
	swg.Path("/hello/world", http.MethodPost, &PathOptions{ID: "createHelloWorld"}).
		PathParams(TestPathParams{}).
		QueryParams(TestQueryParams{}).
		Response(http.StatusOK, TestResponse{}).
		Response(http.StatusNotFound, BaseResponse{}).
		Response(http.StatusInternalServerError, BaseResponse{})

	// Future: Prefix returns path prefix where we can set common Response, PathParams, QueryParams for all Path that
	// are created from it
	// api := swg.Prefix("/api/v1").Response(http.StatusInternalServerError, BaseResponse{})
	// api.Path("user", http.MethodGet)

	b, err := json.Marshal(swg)
	assert.NoError(t, err)
	println(string(b))
}
