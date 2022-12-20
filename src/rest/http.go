package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type RESTInterface struct {
}

func (rest *RESTInterface) GetDevices(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting Devices\n")
}

func (rest *RESTInterface) SetNewDevice(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Set Device\n")
}

func Setup() *mux.Router {
	rest := RESTInterface{}
	router := mux.NewRouter()
	router.HandleFunc("/devices", rest.GetDevices).Methods("GET")
	router.HandleFunc("/devices", rest.SetNewDevice).Methods("POST")
	return router
}

func Start(router *mux.Router) {
	fmt.Println("HTTP serve at 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
