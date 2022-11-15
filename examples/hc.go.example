package main

import (
	"log"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
)

func main() {
	// create an accessory
	info := accessory.Info{Name: "Lamp"}
	ac := accessory.NewSwitch(info)

	// configure the ip transport
	config := hc.Config{Pin: "00102003"}
	t, err := hc.NewIPTransport(config, ac.Accessory)
	if err != nil {
		log.Panic(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()
}
