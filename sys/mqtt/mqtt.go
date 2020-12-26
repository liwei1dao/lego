package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func newSys(options Options) (sys *Mqtt, err error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", options.MqttUrl, options.MqttPort))
	opts.SetClientID(options.ClientID)
	opts.SetUsername(options.UserName)
	opts.SetPassword(options.Password)
	opts.SetDefaultPublishHandler(options.MqttFactory.MessageHandler)
	opts.OnConnect = options.MqttFactory.ConnectHandler
	opts.OnConnectionLost = options.MqttFactory.ConnectionLostHandler
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		err = token.Error()
		return
	}
	sys = &Mqtt{
		client: client,
		token:  token,
	}
	return
}

type Mqtt struct {
	client mqtt.Client
	token  mqtt.Token
}

func (this *Mqtt) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	return this.client.Subscribe(topic, qos, callback)
}

func (this *Mqtt) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	return this.client.Publish(topic, qos, retained, payload)
}
