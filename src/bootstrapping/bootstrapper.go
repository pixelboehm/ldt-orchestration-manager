package bootstrapping

import (
	"encoding/json"
	"log"
	comms "longevity/src/communication"
	mon "longevity/src/ldt-orchestrator/monitor"

	. "longevity/src/database"
	"net/http"
)

type Bootstrapper struct {
	rest        *comms.RESTInterface
	monitor     *mon.Monitor
	waitingList chan Device
}

func NewBootstrapper(monitor *mon.Monitor) *Bootstrapper {
	return &Bootstrapper{
		rest:    comms.NewRestInterface(nil),
		monitor: monitor,
		// waitingList: make(chan Device),
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

	// b.waitingList = make(chan Device)
	// b.waitingList <- device
	// close(b.waitingList)

	name := b.getPariableLDT(device)
	w.Write([]byte(name))
}

func (b *Bootstrapper) getPariableLDT(waiting_device Device) string {
	var pairable_ldt string = b.monitor.GetPairaibleLDT()

	if pairable_ldt == waiting_device.Name {
		return waiting_device.Name
	} else {
		log.Printf("Pairable: %s \t Waiting: %s", pairable_ldt, waiting_device.Name)
	}
	return " "
}
