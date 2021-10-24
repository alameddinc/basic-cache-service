package controllers

import (
	"encoding/json"
	"errors"
	"github.com/alameddinc/ysc/models"
	. "github.com/alameddinc/ysc/schema"
	"github.com/gorilla/mux"
	"net/http"
)

func PostSetValue(w http.ResponseWriter, r *http.Request) {
	requestSchema := SetValueRequestSchema{}
	json.NewDecoder(r.Body).Decode(&requestSchema)
	val, err := models.CreateValue(requestSchema.Key, requestSchema.Value)
	if err != nil {
		jsonFailed(w, err)
	}
	responseSchema := ValueResponseSchema{
		Value:   requestSchema.Value,
		Storage: val.FilenameStamp,
	}
	jsonSuccessfull(w, responseSchema)
}

func GetValueContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	if value, ok := models.CachedValues[key]; ok {
		responseSchema := ValueResponseSchema{
			Value:   value.Content,
			Storage: value.FilenameStamp,
		}
		jsonSuccessfull(w, responseSchema)
		return
	}
	jsonFailed(w, errors.New("Not Found"))
}
