package tags

import "github.com/phonkee/attribs"

var (
	QueryParamParser = attribs.MustNew(QueryParam{})
	PathParamParser  = attribs.MustNew(PathParam{})
)

type QueryParam struct {
	Disabled    bool   `attr:"name=disabled"`
	Name        string `attr:"name=name"`
	Description string `attr:"name=description"`
}

type PathParam struct {
	Disabled    bool   `attr:"name=disabled"`
	Name        string `attr:"name=name"`
	Description string `attr:"name=description"`
}
