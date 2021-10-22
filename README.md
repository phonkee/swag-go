# swag-go

Simple swagger generator to be used for my purposes.

####  Warning
This library is intended to be used in `init` methods, so error handling is basically done with panic.
If you use it differently, please think about this.
The idea is to define all things in init, and then just serve swagger.json (or yaml) which is then cached.

# example

```go
// package level service (so we can access it)
var Service swag.Swagger

func init() {
	Service = swag.New("pat store")
}

type GetPathParams struct {
	ID int `json:"id" swag_description:"primary key in pets database"`
}

type Pet struct {
	ID int `json:"id" swag_description:"unique identifier in database"`
	Name string `json:"name"`
	Born time.Time `json:"born"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreatePetSerializer struct {
	Name string `json:"name" swag_description:"Name of your pet"`
}

type CreatePetValidationError struct {
	Fields map[string]string `json:"fields"`
}

type FilterPetsQuery struct {
	Dogs bool `json:"dogs" swag_description:"only dogs will be returned"`
}

func init() {
	// add get pet by id
    Service.Path("/api/v1/pets/{id}", http.MethodGet).
        // add path params
        PathParams(GetPathParams{}).
        Response(http.StatusOK, Pet{}).
        // in this case 404 does not return any specific response
		Response(http.StatusNotFound).
        Response(http.StatusInternalServerError, ErrorResponse{})
    
    // create new pet endpoint
    Service.Path("/api/v1/pets", http.MethodPost).
        Body(CreatePetSerializer{}).
        Response(http.StatusOK, Pet{}).
        Response(http.StatusBadRequest, CreatePetValidationError{})
    
    // list pets endpoint
    Service.Path("/api/v1/pets", http.MethodGet).
        QueryParams(FilterPetsQuery{}).
        Response(http.StatusOK, []Pet{})
}
```

We have also ability to have shared common properties:

```go
ApiV1 := Service.Prefix("/api/v1/").
	Response(http.StatusNotFound).
	Response(http.StatusInternalServerError, ErrorResponse{})
```

And even this:

```go
type UserIdentifierPathQuery interface {
	ID int `json:"id"`
}

type Order struct {
	ID int `json:"id"`
}

func init() {
    // prepare prefix that identifies user by id
	UsersOrdersApiV1 := Service.Prefix("/api/v1/users/{id}/orders").
		PathParams(UserIdentifierPathQuery{}).
        Response(http.StatusNotFound).
        Response(http.StatusInternalServerError, ErrorResponse{})

	// now get list of orders for user - path will be /api/v1/users/{id}/orders
    UsersOrdersApiV1.Path("", http.MethodGet).
		Response(http.StatusOK, []Order{})

	// return single order by order_id
    UsersOrdersApiV1.Path("{order_id}", http.MethodGet).
		Response(http.StatusOK, Order{})
}

```


# author

Peter Vrba <phonkee@phonkee.eu>
