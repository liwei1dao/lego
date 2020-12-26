package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type DefMqttFactory struct {
}

func (this *DefMqttFactory) ConnectHandler(client mqtt.Client) {
	fmt.Printf("ConnectHandler: succ \n")
}

func (this *DefMqttFactory) ConnectionLostHandler(client mqtt.Client, err error) {
	fmt.Printf("ConnectionLostHandler: succ \n")
}

func (this *DefMqttFactory) MessageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("MessageHandler: topic:%s message:%v \n", msg.Topic(), msg.Payload())
}
