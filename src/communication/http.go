package communication

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
	Run(port int)
	RunWithHTTPS(port int)
	Devices(w http.ResponseWriter, r *http.Request)
	Device(w http.ResponseWriter, r *http.Request)
	CreateDevice(w http.ResponseWriter, r *http.Request)
	UpdateDevice(w http.ResponseWriter, r *http.Request)
	DeleteDevice(w http.ResponseWriter, r *http.Request)
	Database() *sql.DB
	Router() *mux.Router
	SetDatabase(db *sql.DB)
	SetRouter(router *mux.Router)
	CloseDatabase(db *sql.DB)
	initialize()
}

type RESTInterface struct {
	router *mux.Router
	db     *sql.DB
}

func NewRestInterface(db *sql.DB) *RESTInterface {
	return &RESTInterface{
		router: mux.NewRouter(),
		db:     db,
	}
}

func (rest *RESTInterface) Run(port int) {
	rest.initialize()
	log.Printf("HTTP serve at %d\n", port)
	addr := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(addr, rest.router)
	if err != nil {
		panic(err)
	}
}

func (rest *RESTInterface) RunWithHTTPS(port int) {
	rest.initialize()
	log.Printf("HTTPS serve at %d\n", port)
	addr := fmt.Sprintf(":%d", port)
	err := http.ListenAndServeTLS(addr, "server.crt", "server.key", rest.router)
	if err != nil {
		panic(err)
	}
}

func (rest *RESTInterface) Devices(w http.ResponseWriter, r *http.Request) {
	result := []string{}
	respondWithJSON(w, http.StatusOK, result)
}

func (rest *RESTInterface) Device(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid device ID")
		return
	}

	p := Device{ID: id}
	err = p.Device(rest.db)
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
	var device Device
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&device)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid Payload")
		return
	}
	defer r.Body.Close()

	err = device.CreateDevice(rest.db)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, device)
}

func (rest *RESTInterface) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unknown Device ID")
		return
	}

	var device Device
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&device)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid Payload")
		return
	}
	defer r.Body.Close()
	device.ID = id

	err = device.UpdateDevice(rest.db)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, device)
}

func (rest *RESTInterface) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unknown Device ID")
		return
	}

	var device = &Device{ID: id}
	err = device.DeleteDevice(rest.db)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (rest *RESTInterface) Database() *sql.DB {
	return rest.db
}

func (rest *RESTInterface) Router() *mux.Router {
	return rest.router
}

func (rest *RESTInterface) SetDatabase(db *sql.DB) {
	rest.db = db
}

func (rest *RESTInterface) SetRouter(router *mux.Router) {
	rest.router = router
}

func (rest *RESTInterface) CloseDatabase(db *sql.DB) {
	rest.db.Close()
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

func (rest *RESTInterface) initialize() {
	rest.router.HandleFunc("/devices", rest.Devices).Methods("GET")
	rest.router.HandleFunc("/device/{id:[0-9]+}", rest.Device).Methods("GET")
	rest.router.HandleFunc("/device/{id:[0-9]+}", rest.UpdateDevice).Methods("PUT")
	rest.router.HandleFunc("/device/{id:[0-9]+}", rest.DeleteDevice).Methods("DELETE")
	rest.router.HandleFunc("/device", rest.CreateDevice).Methods("POST")
}
