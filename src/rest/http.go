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
	fmt.Printf("Getting Devices")
}

func (rest *RESTInterface) SetNewDevice(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Set Device")
}

func Setup() {
	rest := RESTInterface{}
	router := mux.NewRouter()
	router.HandleFunc("/devices", rest.GetDevices).Methods("GET")
	fmt.Println("HTTP serve at 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
