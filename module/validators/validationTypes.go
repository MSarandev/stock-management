package validators

// InsertStock a custom validation struct for the insertion fields.
type InsertStock struct {
	ID   string `validate:"required,uuid4" json:"id"`
	Name string `validate:"required,alphanumunicode" json:"name"`
}

// UpdateStock a custom validation struct for the update fields.
type UpdateStock struct {
	Name     string `validate:"" json:"name"`
	Quantity int    `validate:"" json:"name"`
}

// GetStock a validator for the single get request.
type GetStock struct {
	ID string `validate:"required,uuid4" json:"id"`
}
