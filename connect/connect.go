package connect

import (
	"fmt"
	"strconv"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

// Will defines the LWT (Last Will & Testament) Message,
// Which will be sent after disconnection.
type Will struct {
	Topic   string
	Message string
	QoS     byte
	Retain  bool
}

// Connection contains the client instance and the connection parameters.
type Connection struct {
	client   *client.Client
	Host     string
	Port     int
	Name     string
	Username string
	Password string
	Will     *Will
}

// New instantiates a connection to the MQTT server
func New(host string, port int, name string, username string, password string) (Connection, error) {
	conn := Connection{Host: host, Port: port, Name: name, Username: username, Password: password}
	err := conn.Connect()

	if err != nil {
		return conn, err
	}

	return conn, nil
}

// Connect to the MQTT server
func (c *Connection) Connect() error {
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
	fmt.Println("[MQTT] Attempting to connect to the broker at", address)

	will := &Will{}

	will.Topic = "axi/status"
	will.Message = "EXIT"

	err := c.client.Connect(&client.ConnectOptions{
		Network:      "tcp",
		Address:      address,
		ClientID:     []byte(c.Name),
		UserName:     []byte(c.Username),
		Password:     []byte(c.Password),
		CleanSession: true,
		WillTopic:    []byte(will.Topic),
		WillMessage:  []byte(will.Message),
		WillQoS:      will.QoS,
		WillRetain:   will.Retain,
	})

	if err != nil {
		fmt.Println("[Error] Connection to the MQTT Broker Failed.")
		return err
	}

	fmt.Println("[MQTT] Connection to the MQTT Broker is Established.")
	return nil
}

// Subscribe to a room
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

// Unsubscribe from the room
func (c *Connection) Unsubscribe(room string) {
	err := c.client.Unsubscribe(&client.UnsubscribeOptions{
		TopicFilters: [][]byte{
			[]byte(room),
		},
	})

	if err != nil {
		panic(err)
	}
}

// On will listen to the topic, and call the handler if new message arises.
func (c *Connection) On(room string, handler func(topic, message string)) {
	c.Subscribe(room, mqtt.QoS0, func(topic, message []byte) {
		handler(string(topic), string(message))
	})
}

// Spy on a room and log the messages as soon as it arrives.
func (c *Connection) Spy(room string) {
	c.On(room, func(topic, message string) {
		fmt.Println("[MQTT: "+topic+"]", message)
	})
}

// Publish will emit a message to that topic.
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

// Send is a shorthand for publishing with QoS of 0.
func (c *Connection) Send(topic, message string) {
	c.Publish(topic, message, mqtt.QoS0, false)
}

// SendRetain will publish a retained message.
func (c *Connection) SendRetain(topic, message string) {
	c.Publish(topic, message, mqtt.QoS0, true)
}

// Close will disconnect from the server.
func (c *Connection) Close() {
	if err := c.client.Disconnect(); err != nil {
		panic(err)
	}

	fmt.Println("[MQTT] Connection Closed; Device is Disconnected.")
}
