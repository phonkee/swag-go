package params

import "github.com/phonkee/tagstruct"

// ParamTag describes parameter struct tag
type ParamTag struct {
	Description string `ts:"name=description"`
	Prefix      string `ts:"name=prefix"`
	Disabled    bool   `ts:"name=disabled"`
}

var (
	// paramTag is parser for param tag struct
	paramTag = tagstruct.New(ParamTag{})
)
