package system

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/phoomparin/axi-device/connect"
)

const (
	device      = "10254"
	statusTopic = "axi/status/#"
)

var isOnline = false
var conn *connect.Connection

// DeviceStat reports the device stat
type DeviceStat struct {
	Health string
	Time   time.Time
}

// Task performs internal tasks
func Task() {
	fmt.Println("[Axi] Status changed to READY.")
	fmt.Println("[Axi] No tasks in queue. Status changed to IDLE, awaiting for tasks.")

	conn.SendRetain("device/status", "ONLINE")

	m := DeviceStat{
		Health: "OK",
		Time:   time.Now(),
	}

	msg, _ := json.Marshal(m)

	conn.SendRetain("device/stat", string(msg))
}

// Exit handles graceful shutdown
func Exit() {
	conn.SendRetain("device/status", "OFFLINE")

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
	connection, err := connect.New("localhost", 1883, "axi-client", "axi", "hello-world")

	if err != nil {
		fmt.Println("[Axi] Connection Failure. Running in OFFLINE mode.")
	} else {
		fmt.Println("[Axi] Connection Established. Running in ONLINE mode.")

		conn = &connection
		isOnline = true
		Listen()
	}
}
