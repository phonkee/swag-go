package swag

type TagPath struct {
	Description string `swag:"description='lorem ipsum dolor sit amet'"`
}

func ParseTagPath(tag string) *TagPath {
	return nil
}
