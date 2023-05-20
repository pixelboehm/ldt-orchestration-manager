package bootstrapping

import (
	"log"
	comms "longevity/src/communication"
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
	w.Write([]byte("hello new device"))
}
