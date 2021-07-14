package kafka

import (
	"fmt"
	"testing"

	"github.com/Shopify/sarama"
)

func Test_sys(t *testing.T) {
	if err := OnInit(map[string]interface{}{
		"StartType": All,
		"Hosts":     []string{"49.235.111.75:9092"},
		"Topics":    []string{"test"},
		"GroupId":   "test-c1",
		"ClientID":  "send",
	}); err != nil {
		fmt.Printf("start sys err:%v", err)
	} else {
		Syncproducer_SendMessage(&sarama.ProducerMessage{
			Topic: "send",
			Value: sarama.StringEncoder("asdasd"),
		})
	}
}
