package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var volume int
var volumeWasSent bool

func getHostname() string {

	hostname, err := os.Hostname()

	if err != nil {
		log.Fatal(err)
	}

	// name.local => name
	firstPart := strings.Split(hostname, ".")[0]

	// maybe we should remove all symbols, but [a-z0-9_-] ?

	return firstPart
}

func getCurrentVolume() int {
	cmd := exec.Command("osascript", "-e", "output volume of (get volume settings)")

	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	stdoutStr := string(stdout)
	stdoutStr = strings.TrimSuffix(stdoutStr, "\n")

	i, err := strconv.Atoi(stdoutStr)
	if err != nil {
		log.Fatal(err)
	}

	return i
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected to MQTT")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Disconnected from MQTT: %v", err)
}

func getMQTTClient() mqtt.Client {
	var broker = "..."
	var port = 1883
	var username = "..."
	var password = "..."

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

func main() {
	var wg sync.WaitGroup

	hostname := getHostname()
	log.Println("Hostname:", hostname)

	mqttClient := getMQTTClient()

	ticker := time.NewTicker(2000 * time.Millisecond)
	wg.Add(1)
	go func() {
		for {
			select {
			case _ = <-ticker.C:
				currentVolume := getCurrentVolume()

				// sending to mqtt only values when they are changed
				if !volumeWasSent || volume != currentVolume {
					token := mqttClient.Publish("mac2mqtt/"+hostname+"/status/volume", 0, false, strconv.Itoa(currentVolume))
					token.Wait()
					volumeWasSent = true
					volume = currentVolume
				}
			}
		}
	}()

	wg.Wait()
	fmt.Println("end")
}
