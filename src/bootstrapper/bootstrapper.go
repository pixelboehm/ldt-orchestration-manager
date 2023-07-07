package bootstrapper

import (
	"encoding/json"
	"fmt"
	"log"
	comms "longevity/src/communication"
	man "longevity/src/ldt-orchestrator/manager"
	mon "longevity/src/monitoring-dependency-manager"
	. "longevity/src/types"
	"net/http"
)

type Bootstrapper struct {
	rest    *comms.RESTInterface
	monitor *mon.Monitor
	manager *man.Manager
}

func NewBootstrapper(monitor *mon.Monitor, manager *man.Manager) *Bootstrapper {
	return &Bootstrapper{
		rest:    comms.NewRestInterface(nil),
		monitor: monitor,
		manager: manager,
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

	ldt_address := b.bootstrap(device)
	w.Write([]byte(ldt_address))
}

func (b *Bootstrapper) bootstrap(waiting_device Device) string {
	result, err := b.monitor.GetLDTAddressForDevice(waiting_device)
	if err != nil {
		log.Println(fmt.Sprint("Bootstrapper: Failed to find pairable LDT", err))
		return " "
	}
	if result != "No pairable LDT available" {
		ldt_address := result
		return ldt_address
	} else {
		b.startSuitableLdt(waiting_device)
		return b.getLDTAddressForDevice(waiting_device)
	}
}

func (b *Bootstrapper) startSuitableLdt(waiting_device Device) {
	b.manager.OptionalScan()
	var full_ldt_specifier string
	for _, ldt := range b.manager.Discovery.SupportedLDTs.LDTs {
		if ldt.Name == waiting_device.Name && ldt.Version[1:] == waiting_device.Version {
			full_ldt_specifier = ldt.User + "/" + ldt.Name + ":" + ldt.Version[1:]
			b.manager.DownloadLDT(full_ldt_specifier)
			break
		}
	}

	process, err := b.manager.RunLDT([]string{"run", full_ldt_specifier})
	if err != nil {
		panic(err)
	}
	b.monitor.Started <- process
}
