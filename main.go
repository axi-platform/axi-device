package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/phoomparin/axi-device/connect"
)

func bootstrap() {
	fmt.Println("[Axi] Bootstrapping Device...")

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	conn := connect.New("iot.eclipse.org", 1883, "Axi Client")
	conn.Listen("hello/#")
	conn.Send("hello/123", "Hello 123! Sending to channel hello/123.")

	conn.Listen("data/#")
	conn.Send("data/5341", "Data 5341! Sending to channel data/5341.")
	conn.Send("hello/25365", "Hello 25365! Sending to channel hello/25365.")

	fmt.Println("[Device] Status changed to READY.")
	fmt.Println("[Device] No tasks in queue. Status changed to IDLE, awaiting for tasks.")

	<-sigc

	conn.Close()
}

func main() {
	bootstrap()
}
