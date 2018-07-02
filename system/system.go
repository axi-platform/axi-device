package system

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/phoomparin/axi-device/connect"
)

const (
	device           = "10254"
	statusTopic      = "device/status/#"
	createQueueTopic = "queue/#/create"
	basePath         = "http://localhost:3030/upload/"
)

var isOnline = false
var conn *connect.Connection

// DeviceStat reports the device stat
type DeviceStat struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}

func reportStatus(status string) {
	m := DeviceStat{
		Status: status,
		Time:   time.Now(),
	}

	msg, _ := json.Marshal(m)

	color.Blue("[Engine] Sending Status Report: %s", string(msg))
	conn.SendRetain("device/"+conn.Name+"/status", string(msg))
}

func download(fileName string) {
	out, err := os.Create(fileName)
	defer out.Close()

	if err != nil {
		fmt.Errorf("Failed to create file: %s", fileName)
	}

	url := basePath + fileName
	res, err := http.Get(url)
	defer res.Body.Close()

	color.Blue("Downloading Document from %s", url)

	if res.StatusCode != http.StatusOK && err != nil {
		fmt.Errorf("Failed to download file: %s %s", url, res.StatusCode)
	}

	io.Copy(out, res.Body)
	color.Green("[PrintAt] Downloaded File: %s", fileName)
}

// Task performs internal tasks
func Task() {
	color.Green("[Engine] Status changed to READY.")
	color.Blue("[Engine] No tasks in queue. Status changed to IDLE, awaiting for tasks.")

	reportStatus("ONLINE")
}

// Exit handles graceful shutdown
func Exit() {
	reportStatus("OFFLINE")

	if isOnline {
		conn.Close()
	}

	fmt.Println("[Engine] Shutting Down... Bye!")
}

// QueueData is a queue data
type QueueData struct {
	Files []string `json:"files"`
}

func handleQueue(room, msg string) {
	fmt.Println("Incoming Queue:", room, msg)

	data := QueueData{}
	json.Unmarshal([]byte(msg), &data)

	fmt.Println("Queue:", data)

	for _, file := range data.Files {
		download(file)
	}

	parts := strings.Split(room, "/")[:4]
	completedRoom := strings.Join(parts, "/") + "/completed"

	conn.SendRetain(completedRoom, "1")
}

func Listen() {
	conn.On(statusTopic, func(room, msg string) {
		fmt.Println("Status Message:", msg)
	})

	conn.On(createQueueTopic, handleQueue)

	conn.Spy("axi/#")
	conn.Spy("hello/#")
}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)

	if ok {
		return value
	}

	return fallback
}

func Boot() {
	yellow := color.New(color.FgYellow, color.Bold)
	green := color.New(color.FgGreen, color.Bold)

	yellow.Println("--- Axi Engine v0.2 | Axi Platform ---")
	color.Blue("[Engine] Bootstrapping Device...")

	port := 1883
	endpoint := getEnv("AXI_ENDPOINT", "localhost")
	deviceID := getEnv("AXI_ID", "printat-demo")
	deviceSecret := getEnv("AXI_SECRET", "printat-demo")

	color.Blue("[Engine] Authenticating with Axi Platform as %s...", deviceID)
	connection, err := connect.New(endpoint, port, deviceID, deviceID, deviceSecret)

	if err != nil {
		color.Red("[Engine] Connection Failure. Running in OFFLINE mode.")
	} else {
		green.Println("[Engine] Connection Established! Running in ONLINE mode.")

		conn = &connection
		isOnline = true
		Listen()
	}
}
