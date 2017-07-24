package system

import (
	"fmt"

	"github.com/phoomparin/axi-device/connect"
)

const (
	device      = "10254"
	statusTopic = "axi/status/#"
)

var isOnline = false
var conn *connect.Connection

func Task() {
	fmt.Println("[Axi] Status changed to READY.")
	fmt.Println("[Axi] No tasks in queue. Status changed to IDLE, awaiting for tasks.")
	// conn.SendRetain("axi/status/A", "ONLINE")
}

func Exit() {
	if isOnline {
		conn.Close()
	}

	fmt.Println("[Axi] Shutting Down... Bye!")
}

func Listen() {
	conn.On(statusTopic, func(room, msg string) {
		fmt.Printf("Status Message from %s: %s\n", room[len(statusTopic)-1:], msg)
	})

	conn.Spy("axi/#")
	conn.Spy("hello/#")
}

func Boot() {
	fmt.Println("[Axi] Bootstrapping Device...")
	connection, err := connect.New("localhost", 1883, "Axi Client", "", "")

	if err != nil {
		fmt.Println("[Axi] Connection Failure. Running in OFFLINE mode.")
	} else {
		fmt.Println("[Axi] Connection Established. Running in ONLINE mode.")

		conn = &connection
		isOnline = true
		Listen()
	}
}
