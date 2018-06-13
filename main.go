package main

import (
	"os"
	"os/signal"

	"github.com/phoomparin/axi-device/system"
)

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	system.Boot()
	system.Task()

	system.OfflineHandler()

	// beacon.Advertise("phoom.in.th")
	// beacon.Broadcast("com.xyz", "helloworld")

	<-sigc

	defer system.Exit()
}
