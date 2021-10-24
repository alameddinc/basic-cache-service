package models

type ValueHistory struct {
	Value *Value
	Op    string `json:"op"`
}

var ValueHistories map[string]map[string]*ValueHistory

func init() {
	ValueHistories = make(map[string]map[string]*ValueHistory)
}
