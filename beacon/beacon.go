package beacon

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Telemetry creates a Eddystone-TLM Beacon
// tlm(battery, temp in celsius, count of advertise frames, times since reboot, Tx Power)
func Telemetry() {
	tlm := NewEddystoneTLMBeacon(500, 22.0, 100, 1000, -20)
	tlm.Advertise()
}

func hash(input string, length int) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	text := strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))[:length]

	return text
}

// Broadcast starts an Eddystone-UID Beacon
func Broadcast(domain, identifier string) {
	var namespace = hash(domain, 20)
	var instance = hash(identifier, 12)

	exec.Command("sh", "-c", "'sudo hciconfig hci0 down'")

	fmt.Printf("[UID Beacon] Attempting to advertise as %s@%s.\n", identifier, domain)
	fmt.Printf("[UID Beacon] Namespace is %s and Instance is %s.\n", namespace, instance)

	uid := NewEddystoneUIDBeacon(namespace, instance, -20)
	err := uid.Advertise()

	if err != nil {
		fmt.Println("[Beacon Error] Cannot Broadcast Beacon.", err)
		fmt.Println(`
Before Starting: sudo hciconfig hci0 down
Please run this binary as root, or sudo setcap 'cap_net_raw,cap_net_admin=eip' <binary>
		`)
		defer os.Exit(127)
		return
	}

	fmt.Printf("[UID Beacon] Broadcasting Eddystone-UID ..")
}

// Advertise starts an Eddystone-URL
func Advertise(url string) {
	edb := NewEddystoneURLBeacon(url, -20)
	edb.Advertise()
}

// BroadcastIBeacon creates an IBeacon
func BroadcastIBeacon(uuid string, name string) {
	ibc := NewIBeacon(uuid, name, -20)
	ibc.AddBatteryService()
	ibc.AddCountService()
	ibc.SetiBeaconVersion(1, 1)
	ibc.Advertise()
}
