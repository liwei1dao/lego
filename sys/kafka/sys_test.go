package kafka

import (
	"fmt"
	"testing"
	"time"

	"github.com/Shopify/sarama"
)

func Test_sys(t *testing.T) {
	if err := OnInit(map[string]interface{}{
		"StartType":                 Asyncproducer,
		"Hosts":                     []string{"172.20.27.126:9092", "172.20.27.127:9092", " 172.20.27.128:9092"},
		"Topics":                    []string{"test"},
		"GroupId":                   "test-c1",
		"ClientID":                  "send",
		"Net_KeepAlive":             0,
		"Producer_MaxMessageBytes":  201231293891,
		"Net_DialTimeout":           5 * time.Second,
		"Producer_Compression":      sarama.CompressionGZIP,
		"Producer_CompressionLevel": 1,
		"Producer_Retry_Max":        3,
	}); err != nil {
		fmt.Printf("start sys err:%v", err)
	} else {
		fmt.Printf("start sys succ")
		// Syncproducer_SendMessage(&sarama.ProducerMessage{
		// 	Topic: "send",
		// 	Value: sarama.StringEncoder("asdasd"),
		// })
	}
}
