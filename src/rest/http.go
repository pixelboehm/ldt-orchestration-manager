package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func GetDevices(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting Devices")
}

func SetNewDevice(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Set Device")
}

func Setup() {
	router := mux.NewRouter()
	router.HandleFunc("/devices", GetDevices).Methods("GET")
	fmt.Println("HTTP serve at 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
