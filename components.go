package swag

func newComponents() Components {
	return &components{
		schemas: newSchemas(),
	}
}

type components struct {
	schemas Schemas
}

func (c *components) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func (c *components) GetSchema(i interface{}) string {
	return c.schemas.GetRef(i)
}
