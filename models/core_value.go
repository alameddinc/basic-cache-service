package models

type CoreValue struct {
	Key     string `json:"key"`
	Content string `json:"content"`
}

func (v CoreValue) CreateBlankValue() *Value {
	return &Value{v, ""}
}
