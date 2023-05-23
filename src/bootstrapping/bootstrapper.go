package bootstrapping

import (
	"encoding/json"
	"log"
	comms "longevity/src/communication"

	. "longevity/src/database"
	"net/http"
)

type Bootstrapper struct {
	rest        *comms.RESTInterface
	waitingList chan Device
}

func NewBootstrapper() *Bootstrapper {
	return &Bootstrapper{
		rest:        comms.NewRestInterface(nil),
		waitingList: make(chan Device),
	}
}

func (b *Bootstrapper) Run(port int) {
	b.rest.AddCustomHandler("/register", b.registration)
	b.rest.Run(port)
}

func (b *Bootstrapper) registration(w http.ResponseWriter, r *http.Request) {
	log.Println("Bootstrapper: A new registration request came")

	var device Device
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&device); err != nil {
		log.Println("Bootstrapper: Decoding Error: ", err)
	}
	defer r.Body.Close()
	b.waitingList <- device
	w.Write([]byte("192.168.188.56"))
}
