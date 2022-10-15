package main

import "github.com/yuin/stagparser"

type User struct {
	Name string `swag:"description='lorem ipsum dolor sit amet' example='bob'"`
}

func main() {
	user := &User{"bob"}
	definitions, err := stagparser.ParseStruct(user, "swag")
	_, _ = definitions, err
}
