package bootstrapping

import (
	"encoding/json"
	"fmt"
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

	ldt_address := b.getLDTAddressForDevice(device)
	w.Write([]byte(ldt_address))
}

func (b *Bootstrapper) getLDTAddressForDevice(waiting_device Device) string {
	ldtAddress, err := b.monitor.GetLDTAddressForDevice(waiting_device)
	if err != nil {
		log.Println(fmt.Sprint("Bootstrapper: Failed to find pairable LDT", err))
		return " "
	}

	return ldtAddress
}
