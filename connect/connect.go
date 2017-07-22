package connect

import (
	"fmt"
	"strconv"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

// Connection defines the conn obj
type Connection struct {
	client *client.Client
	Host   string
	Port   int
	Name   string
}

// Connect to the MQTT server
func (c *Connection) Connect() {
	c.client = client.New(&client.Options{
		ErrorHandler: func(err error) {
			fmt.Println("[MQTT Error]", err)
		},
	})
	defer c.client.Terminate()

	if c.Port == 0 {
		c.Port = 1883
	}

	if c.Name == "" {
		c.Name = "Unnamed Axi Client"
	}

	address := string(c.Host) + ":" + strconv.Itoa(c.Port)
	fmt.Println("[MQTT] Attempting to connect to", address)

	err := c.client.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  address,
		ClientID: []byte(c.Name),
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("[Axi] Connection Established. Running in ONLINE mode.")
}

func (c *Connection) Subscribe(room string, qos byte, handler func(topic, message []byte)) {
	err := c.client.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte(room),
				QoS:         qos,
				Handler:     handler,
			},
		},
	})

	if err != nil {
		panic(err)
	}
}

func (c *Connection) On(room string, handler func(topic, message string)) {
	c.Subscribe(room, mqtt.QoS0, func(topic, message []byte) {
		handler(string(topic), string(message))
	})
}

// Listen to some rooms
func (c *Connection) Listen(room string) {
	c.On(room, func(topic, message string) {
		fmt.Println("[MQTT: "+topic+"]", message)
	})
}

func (c *Connection) Publish(topic string, message string, qos byte, retain bool) {
	err := c.client.Publish(&client.PublishOptions{
		QoS:       qos,
		Retain:    retain,
		TopicName: []byte(topic),
		Message:   []byte(message),
	})

	if err != nil {
		panic(err)
	}
}

func (c *Connection) Send(topic, message string) {
	c.Publish(topic, message, mqtt.QoS0, false)
}

func (c *Connection) SendRetain(topic, message string) {
	c.Publish(topic, message, mqtt.QoS0, true)
}

func (c *Connection) Close() {
	if err := c.client.Disconnect(); err != nil {
		panic(err)
	}

	fmt.Println("[MQTT] Connection Closed.")
}
