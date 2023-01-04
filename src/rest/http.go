package rest

import (
	"encoding/json"
	"fmt"
	"log"
	d "longevity/src/model"
	"net/http"

	"github.com/gorilla/mux"
)

type RESTInterface struct {
	Router *mux.Router
}

type JsonResponse struct {
	Type string
	Data string
}

func NewRestInterface() *RESTInterface {
	return &RESTInterface{
		Router: mux.NewRouter(),
	}
}

func (rest *RESTInterface) GetDevices(w http.ResponseWriter, r *http.Request) {
	response := JsonResponse{Type: "success", Data: "foo"}

	json.NewEncoder(w).Encode(response)
}

func (rest *RESTInterface) SetNewDevice(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	macAddress := r.FormValue("macAddress")
	twin := r.FormValue("twin")
	version := r.FormValue("version")

	if containsEmptyString(name, macAddress) == false {
		d.NewDevice(name, macAddress, twin, version)
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

func containsEmptyString(formValues ...string) bool {
	for _, s := range formValues {
		if s == "" {
			return true
		}
	}
	return false
}
