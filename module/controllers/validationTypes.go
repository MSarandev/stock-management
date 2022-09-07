package controllers

type InsertStock struct {
	ID   string `validate:"required,uuid4" json:"id"`
	Name string `validate:"required,alphanumunicode,len=1" json:"name"`
}

type UpdateStock struct {
	Name     string `validate:"" json:"name"`
	Quantity int    `validate:"" json:"name"`
}
