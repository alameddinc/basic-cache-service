package models

import (
	"errors"
	"fmt"
	"time"
)

var CachedValues map[string]*Value
var NextTimestamp int64

type Value struct {
	CoreValue
	FilenameStamp string `json:"filename_stamp"`
}

func init() {
	CachedValues = make(map[string]*Value)
	NextTimestamp = time.Now().Unix()
}

func CreateValue(key string, content string) (*Value, error) {
	v, err := FetchValue(key)
	if err != nil {
		v = CoreValue{Key: key}.CreateBlankValue()
	}
	return v.Set(content)
}

func FetchValue(key string) (*Value, error) {
	if v, ok := CachedValues[key]; ok {
		return v, nil
	}
	return nil, errors.New("Not Found")
}

func (v *Value) Set(content string) (*Value, error) {
	// Cache on File
	if v.FilenameStamp == "" {
		v.FilenameStamp = fmt.Sprintf("%d.txt", NextTimestamp)
	}
	if _, ok := ValueHistories[v.FilenameStamp]; !ok {
		ValueHistories[v.FilenameStamp] = map[string]*ValueHistory{}
	}
	ValueHistories[v.FilenameStamp][v.Key] = &ValueHistory{
		Value: v,
		Op:    "w",
	}
	// Set on Cache
	v.Content = content
	CachedValues[v.Key] = v
	return v, nil
}

func (v Value) Get() string {
	return v.Content
}

func (v *Value) Delete() error {
	//delete on File
	if _, ok := ValueHistories[v.FilenameStamp]; !ok {
		ValueHistories[v.FilenameStamp] = map[string]*ValueHistory{}
	}
	ValueHistories[v.FilenameStamp][v.Key].Op = "d"
	// delete on Cache
	delete(CachedValues, v.Key)
	return errors.New("not found")
}
