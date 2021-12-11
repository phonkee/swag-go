# swag-go

Simple swagger generator to be used for my purposes.

####  Warning
This library is intended to be used in `init` methods, so error handling is basically done with panic.
If you use it differently, please think about this.
The idea is to define all things in init, and then just serve swagger.json (or yaml) which is then cached.

# example

```go
// package level Service (so we can access it) from handlers
// usually put in domain (api) package so it's accessible everywhere
var Service swag.Swagger

func init() {
	// initialize swag Swagger
	Service = swag.New("pet store")
}

type GetPetPathParams struct {
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

type FieldValidationError struct {
	Fields map[string]string `json:"fields"`
}

type FilterPetsQuery struct {
	Dogs bool `json:"dogs" swag_description:"only dogs will be returned"`
}

func init() {
	// get single pet by id
    Service.Path("/api/v1/pets/{id}", http.MethodGet).
        // add path params
        PathParams(GetPetPathParams{}).
        Response(http.StatusOK, Pet{}).
        // in this case 404 does not return any specific response
		Response(http.StatusNotFound, nil).
        Response(http.StatusInternalServerError, ErrorResponse{})
    
    // create new pet endpoint
    Service.Path("/api/v1/pets", http.MethodPost).
        Body(CreatePetSerializer{}).
        Response(http.StatusOK, Pet{}).
        Response(http.StatusBadRequest, FieldValidationError{})
    
    // list pets endpoint
    Service.Path("/api/v1/pets", http.MethodGet).
        QueryParams(FilterPetsQuery{}).
        Response(http.StatusOK, []Pet{})
}
```

# prefix

We have also ability to have shared common properties by creating prefixes.
Prefixes share Responses, Path Params, Query Params and of course path prefix.

```go
ApiV1 := Service.Prefix("/api/v1/").
	Response(http.StatusNotFound, nil).
    Response(http.StatusUnauthorized, nil).
	Response(http.StatusInternalServerError, ErrorResponse{})
```

And even this more complicated example works:

```go
type UserIdentifierPathQuery struct {
	ID int `json:"id"`
}

type Order struct {
	ID int `json:"id"`
}

type OrderCacheQueryParams struct {
	NoCache bool `json:"no_cache" swag_description:"when true, orders will be fetched from database"`
}

func init() {
    // prepare prefix that identifies user by id
	UsersOrdersApiV1 := Service.Prefix("/api/v1/users/{id}/orders").
		// path params will be inherited in all paths derived from this prefix
		PathParams(UserIdentifierPathQuery{}).
		// query params will be inherited in all paths derived from this prefix 
		QueryParams(OrderCacheQueryParams{}).
        // responses will be inherited in all paths derived from this prefix
		Response(http.StatusNotFound, nil).
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
