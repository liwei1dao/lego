package kafka_test

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/liwei1dao/lego/sys/kafka"
)

func Test_sys(t *testing.T) {
	if err := kafka.OnInit(map[string]interface{}{
		"StartType":                 kafka.Asyncproducer,
		"Hosts":                     []string{"172.20.27.126:9092", "172.20.27.127:9092", "172.20.27.128:9092"},
		"Topics":                    []string{"ETL-IN-CX20211013641233861142316500871138900239"},
		"GroupId":                   "liwei3dao",
		"ClientID":                  "test",
		"Net_KeepAlive":             0,
		"Producer_MaxMessageBytes":  201231293891,
		"Net_DialTimeout":           5 * time.Second,
		"Producer_Compression":      sarama.CompressionGZIP,
		"Producer_CompressionLevel": 1,
		"Producer_Retry_Max":        3,
		"Sasl_Enable":               true,
		"Sasl_Mechanism":            sarama.SASLTypeGSSAPI,
	}, kafka.SetSasl_GSSAPI(sarama.GSSAPIConfig{
		AuthType:           sarama.KRB5_KEYTAB_AUTH,
		Realm:              "TUJL.COM",
		ServiceName:        "kafka",
		Username:           "kafka/sjzt-wuhan-13",
		KeyTabPath:         "./sjzt-wuhan-13.keytab",
		KerberosConfigPath: "./krb5.conf",
		DisablePAFXFAST:    true,
	})); err != nil {
		fmt.Printf("start sys err:%v", err)
	} else {
		fmt.Printf("start sys succ")
	}
}

func Test_ConsumerGroup(t *testing.T) {
	if err := kafka.OnInit(map[string]interface{}{
		"StartType":                 kafka.Consumer,
		"Hosts":                     []string{"172.20.27.98:39201"},
		"Topics":                    []string{"ETL-IN-CX20211013641233861142316500871138900249"},
		"GroupId":                   "liwei4dao",
		"ClientID":                  "test",
		"Net_KeepAlive":             0,
		"Producer_MaxMessageBytes":  201231293891,
		"Net_DialTimeout":           5 * time.Second,
		"Producer_Compression":      sarama.CompressionGZIP,
		"Producer_CompressionLevel": 1,
		"Producer_Retry_Max":        3,
	}); err != nil {
		fmt.Printf("start sys err:%v\n", err)
	} else {
		fmt.Printf("start sys succ\n")
		go func() {
			for v := range kafka.Consumer_Messages() {
				fmt.Printf("Consumer_Messages:%d\n", v.Partition)
			}
		}()
	}
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		fmt.Printf("terminating: via signal\n")
	}
	if err := kafka.Consumer_Close(); err != nil {
		fmt.Printf("Consumer_Close err: %v\n", err)
		return
	} else {
		fmt.Printf("Consumer_Close end: %v\n", err)
	}
}
