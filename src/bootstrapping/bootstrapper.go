package bootstrapping

import (
	"encoding/json"
	"fmt"
	"log"
	comms "longevity/src/communication"

	. "longevity/src/database"
	"net/http"
)

type Bootstrapper struct {
	rest *comms.RESTInterface
}

func NewBootstrapper() *Bootstrapper {
	return &Bootstrapper{
		rest: comms.NewRestInterface(nil),
	}
}

func (b *Bootstrapper) Run(port int) {
	b.rest.AddCustomHandler("/register", registration)
	b.rest.Run(port)
}

func registration(w http.ResponseWriter, r *http.Request) {
	log.Println("Bootstrapper: A new registration request came")

	var device Device
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&device); err != nil {
		log.Println("Bootstrapper: Decoding Error: ", err)
	}
	defer r.Body.Close()
	var result string = fmt.Sprintf("New Device: %s\t%s\t%s\n", device.Name, device.MacAddress, device.Version)
	w.Write([]byte(result))
}
