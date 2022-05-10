package mqtt

import mqtt "github.com/eclipse/paho.mqtt.golang"

type (
	ISys interface {
		Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token
		Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token
	}
	IMqttFactory interface {
		ConnectHandler(client mqtt.Client)
		ConnectionLostHandler(client mqtt.Client, err error)
		MessageHandler(client mqtt.Client, msg mqtt.Message)
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}
func NewSys(option ...Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

func Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	return defsys.Subscribe(topic, qos, callback)
}

func Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	return defsys.Publish(topic, qos, retained, payload)
}
