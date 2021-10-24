package routes

import (
	"fmt"
	. "github.com/alameddinc/ysc/controllers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var mainRouter *mux.Router

func init() {
	mainRouter = mux.NewRouter()
}

func createSubRouter(path string) *mux.Router {
	return mainRouter.PathPrefix(fmt.Sprintf("/%s", path)).Subrouter()
}

func Handler() {
	Routes()
	log.Fatal(http.ListenAndServe(":8080", mainRouter))
}

func Routes() {
	storageRoutes()
}

func storageRoutes() {
	subRouter := createSubRouter("storage")
	subRouter.HandleFunc("/set", PostSetValue).Methods("POST")
	subRouter.HandleFunc("/get/{key}", GetValueContent).Methods("GET")
}
