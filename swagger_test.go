package swag

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestBody struct {
	Something string `json:"something"`
}

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
		// add path params
		PathParams(TestPathParams{}).
		// add query params
		QueryParams(TestQueryParams{}).
		// add body definition
		Body(TestBody{}).
		// add responses
		Response(http.StatusTeapot, nil).
		Response(http.StatusOK, TestResponse{}).
		Response(http.StatusNotFound, BaseResponse{}).
		Response(http.StatusInternalServerError, BaseResponse{})

	b, err := json.Marshal(swg)
	assert.NoError(t, err)
	_ = b
}
