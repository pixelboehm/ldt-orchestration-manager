package rest

import (
	"fmt"
	"log"
	d "longevity/src/model"
	"net/http"

	"github.com/gorilla/mux"
)

type RESTInterface struct {
	Router *mux.Router
}

func NewRestInterface() *RESTInterface {
	return &RESTInterface{
		Router: mux.NewRouter(),
	}
}

func (rest *RESTInterface) GetDevices(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting Devices\n")
}

func (rest *RESTInterface) SetNewDevice(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name != "" {
		d.NewDevice(name, "", "", "")
		fmt.Printf("New Device with name %s created.\n", name)
	} else {
		http.Error(w, "No Name Specified for Device", http.StatusBadRequest)
	}
}

func (rest *RESTInterface) Setup() {
	rest.Router.HandleFunc("/devices/", rest.GetDevices).Methods("GET")
	rest.Router.HandleFunc("/devices/", rest.SetNewDevice).Methods("POST")
}

func (rest *RESTInterface) Start() {
	fmt.Println("HTTP serve at 8000")
	log.Fatal(http.ListenAndServe(":8000", rest.Router))
}
