# swag-go

Simple swagger package written for my purposes

# example

```go
// create service (swagger)
service := swag.New("pat store")

// add post method
service.Path("/hello/world", http.MethodPost).
    // add path params
    PathParams(PathParams{}).
    // add query params
    QueryParams(QueryParams{}).
    // add body definition
    Body(Body{}).
    // add responses
    Response(http.StatusOK, OkResponse{}).
    Response(http.StatusNotFound, BaseResponse{}).
    Response(http.StatusInternalServerError, BaseResponse{})

```


# author

Peter Vrba <phonkee@phonkee.eu>
