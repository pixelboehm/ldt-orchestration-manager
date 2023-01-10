package rest

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	. "longevity/src/database"
	"net/http"
	"strconv"

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
	respondWithJSON(w, http.StatusOK, result)
}

func (rest *RESTInterface) GetDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid device ID")
		return
	}

	p := DB_Device{ID: id}
	err = p.GetDevice(rest.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "device not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (rest *RESTInterface) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var device DB_Device
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&device)

	if err != nil {
		log.Fatal(err)
		respondWithError(w, http.StatusInternalServerError, "Invalid Payload")
		return
	}
	defer r.Body.Close()

	err = device.CreateDevice(rest.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, device)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response, _ := json.Marshal(payload)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

func (rest *RESTInterface) setup() {
	rest.Router.HandleFunc("/devices", rest.GetDevices).Methods("GET")
	rest.Router.HandleFunc("/device/{id:[0-9]+}", rest.GetDevice).Methods("GET")
	rest.Router.HandleFunc("/device", rest.CreateDevice).Methods("POST")
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
