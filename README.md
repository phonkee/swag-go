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

type Pet struct {
	ID int `json:"id"`
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
    Response(http.StatusOK, Pet{}).
    Response(http.StatusNotFound).
    Response(http.StatusInternalServerError, ErrorResponse{})

type CreatePetSerializer struct {
	Name string `json:"name" swag_description:"Name of your pet"`
}

type CreatePetValidationError struct {
	Fields map[string]string `json:"fields"`
}

// Now create new pet endpoint
service.Path("/api/v1/pets", http.MethodPost).
	Body(CreatePetSerializer{}).
	Response(http.StatusOK, Pet{}).
	Response(http.StatusBadRequest, CreatePetValidationError{})


type FilterPetsQuery struct {
	Dogs bool `json:"dogs" swag_description:"only dogs will be returned"`
}

// now list endpoint
service.Path("/api/v1/pets", http.MethodGet).
    QueryParams(FilterPetsQuery{}).
    Response(http.StatusOK, []Pet{})
```


# author

Peter Vrba <phonkee@phonkee.eu>
