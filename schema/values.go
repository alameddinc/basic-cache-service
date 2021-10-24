package schema

type SetValueRequestSchema struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ValueResponseSchema struct {
	Value   string `json:"value"`
	Storage string `json:"storage"`
}
