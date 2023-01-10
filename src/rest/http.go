package rest

import (
	"database/sql"
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
	DB     *sql.DB
}

type JsonResponse struct {
	Type string
	Data string
}

func NewRestInterface(db *sql.DB) *RESTInterface {
	return &RESTInterface{
		Router: mux.NewRouter(),
		DB:     db,
	}
}

func (rest *RESTInterface) GetDevices(w http.ResponseWriter, r *http.Request) {
	result := []string{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(result)
	w.Write(response)
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
	rest.Router.HandleFunc("/devices", rest.GetDevices).Methods("GET")
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
