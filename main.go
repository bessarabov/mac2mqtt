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

var hostname string

func getHostname() string {

	hostname, err := os.Hostname()

	if err != nil {
		log.Fatal(err)
	}

	// "name.local" => "name"
	firstPart := strings.Split(hostname, ".")[0]

	// maybe we should remove all symbols, but [a-z0-9_-] ?

	return firstPart
}

func getMuteStatus() bool {
	cmd := exec.Command("osascript", "-e", "output muted of (get volume settings)")

	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	stdoutStr := string(stdout)
	stdoutStr = strings.TrimSuffix(stdoutStr, "\n")

	b, err := strconv.ParseBool(stdoutStr)
	if err != nil {
		log.Fatal(err)
	}

	return b
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

func setVolume(i int) {

	cmd := exec.Command("osascript", "-e", "set volume output volume "+strconv.Itoa(i))

	_, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

}

// true - turn mute on
// false - turn mute off
func setMute(b bool) {

	cmd := exec.Command("osascript", "-e", "set volume output muted "+strconv.FormatBool(b))

	_, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

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

	opts.SetWill(getTopicPrefix()+"/status/alive", "false", 0, false)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	token := client.Publish(getTopicPrefix()+"/status/alive", 0, false, "true")
	token.Wait()

	log.Println("Sending 'true' to topic: " + getTopicPrefix() + "/status/alive")

	return client
}

func getTopicPrefix() string {
	return "mac2mqtt/" + hostname
}

func listen(client mqtt.Client, topic string) {
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {

		if msg.Topic() == getTopicPrefix()+"/command/volume" {

			i, err := strconv.Atoi(string(msg.Payload()))
			if err == nil {
				setVolume(i)

				updateVolume(client)
				updateMute(client)

			} else {
				log.Println("Incorrect value")
			}

		}

		if msg.Topic() == getTopicPrefix()+"/command/mute" {

			b, err := strconv.ParseBool(string(msg.Payload()))
			if err == nil {
				setMute(b)

				updateVolume(client)
				updateMute(client)

			} else {
				log.Println("Incorrect value")
			}

		}

	})
}

func updateVolume(client mqtt.Client) {
	token := client.Publish(getTopicPrefix()+"/status/volume", 0, false, strconv.Itoa(getCurrentVolume()))
	token.Wait()
}

func updateMute(client mqtt.Client) {
	token := client.Publish(getTopicPrefix()+"/status/mute", 0, false, strconv.FormatBool(getMuteStatus()))
	token.Wait()
}

func main() {

	var wg sync.WaitGroup

	hostname = getHostname()
	log.Println("Hostname:", hostname)

	mqttClient := getMQTTClient()

	volumeTicker := time.NewTicker(2 * time.Second)

	go listen(mqttClient, getTopicPrefix()+"/command/#")

	wg.Add(1)
	go func() {
		for {
			select {
			case _ = <-volumeTicker.C:
				updateVolume(mqttClient)
				updateMute(mqttClient)
			}
		}
	}()

	wg.Wait()
	fmt.Println("end")
}
