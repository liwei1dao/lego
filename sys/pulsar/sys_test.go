package pulsar_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	lgpulsar "github.com/liwei1dao/lego/sys/pulsar"
)

func Test_Sys(t *testing.T) {
	if err := lgpulsar.OnInit(map[string]interface{}{
		"StartType": lgpulsar.Producer,
		"PulsarUrl": "pulsar://172.20.27.145:6650",
		"Topics":    []string{"liwei1dao_test"},
	}); err != nil {
		fmt.Printf("start sys err:%v", err)
	} else {
		fmt.Printf("start sys succ")
		lgpulsar.Producer_SendAsync() <- &pulsar.ProducerMessage{Payload: []byte("hello word")}
		time.Sleep(time.Second * 2)
	}
}
