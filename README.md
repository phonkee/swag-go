# swag-go

Simple swagger package written for my purposes

# example

```go
// create service (swagger)
service := swag.New("pat store")

// now we create some structs that describe endpoint

type GetPathParams struct {
	ID int `json:"id"`
}

type PetResponse struct {
	Name string `json:"name"`
	Born time.Time `json:"born"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// add post method
service.Path("/api/v1/pets/{id}", http.MethodGet).
    // add path params
    PathParams(GetPathParams{}).
    // add query params
    //QueryParams(QueryParams{}).
    // add body definition
    //Body(Body{}).
    // add responses
    Response(http.StatusOK, PetResponse{}).
    Response(http.StatusNotFound).
    Response(http.StatusInternalServerError, ErrorResponse{})

```


# author

Peter Vrba <phonkee@phonkee.eu>
