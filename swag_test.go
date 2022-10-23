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
		Description: "Pet store swag implementation",
		Version:     "1.0.0",
		License: &License{
			Name: "MIT",
		},
	})
	assert.NotNil(t, swg)

	// add post method
	swg.Path("/hello/world", http.MethodPost, &PathOptions{ID: "createHelloWorld"}).
		// add pathImpl pathParams
		PathParams(TestPathParams{}).
		// add query pathParams
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

type TestParamsParams struct {
	Str string
}

func TestParams(t *testing.T) {
	t.Run("simple pathParams", func(t *testing.T) {
		swg := New("hello")
		swg.Path("/hello", http.MethodGet).
			PathParams(TestParamsParams{})
	})
}

func TestPrefix(t *testing.T) {
	t.Run("test prefix", func(t *testing.T) {
		swg := New("hello")
		swgPrefix := swg.Prefix("/api/v1").
			Response(http.StatusTeapot, nil)
		swgUserPrefix := swgPrefix.Prefix("user")
		p := swgUserPrefix.Path("", http.MethodGet).
			Response(http.StatusNotFound, nil).
			Response(http.StatusUnauthorized, nil)

		// get specs
		spe := p.(*pathImpl).spec()
		assert.Equal(t, 1, len(spe.Paths))
		assert.NotPanics(t, func() {
			_ = spe.Paths["/api/v1/user"]
		})

		data, err := json.MarshalIndent(swg, "  ", "  ")
		assert.NoError(t, err)
		println(string(data))

		//swg.Debug()
	})
}
