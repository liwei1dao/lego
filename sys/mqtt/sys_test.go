package mqtt

import (
	"fmt"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Test_sys(t *testing.T) {
	// sub(client)
	if err := OnInit(map[string]interface{}{
		"MqttUrl":  "127.0.0.1",
		"MqttPort": 1883,
		"ClientID": "go_mqtt_client",
		"UserName": "emqx",
		"Password": "public",
	}); err != nil {
		fmt.Printf("start sys err:%v", err)
	} else {
		sub()
		publish()
	}
}

func publish() {
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := Publish("topic/test", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}

func sub() {
	token := Subscribe("topic/test", 1, func(c mqtt.Client, m mqtt.Message) {
		fmt.Printf("topic/test %d \n", m.Payload())
	})
	token.Wait()
	fmt.Printf("Subscribed to topic: %s \n", "topic/test")
}
