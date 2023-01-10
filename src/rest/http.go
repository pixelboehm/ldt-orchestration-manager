package rest

import (
	"encoding/json"
	"fmt"
	"log"
	d "longevity/src/model"
	"net/http"

	"github.com/gorilla/mux"
)

type API interface {
	GetDevices(w http.ResponseWriter, r *http.Request)
	SetNewDevice(w http.ResponseWriter, r *http.Request)
	Start()
	setup()
}

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

func (rest *RESTInterface) setup() {
}

func (rest *RESTInterface) Start() {
	rest.setup()
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
