package controllers

import (
	"encoding/json"
	"net/http"
)

func jsonSuccessfull(w http.ResponseWriter, value interface{}) {
	json.NewEncoder(w).Encode(value)
}

func jsonFailed(w http.ResponseWriter, error error) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(error)
}
