package rabbitmq_test

import (
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/sys/rabbitmq"
)

func Test_Sys_Producer(t *testing.T) {
	if err := rabbitmq.OnInit(map[string]interface{}{
		"RabbitmqUrl": "amqp://root:li13451234@localhost:5672/",
	}); err != nil {
		fmt.Printf("start sys err:%v", err)
	} else {

	}

}
func Test_Sys_Consumer(t *testing.T) {
	if err := rabbitmq.OnInit(map[string]interface{}{
		"RabbitmqUrl": "amqp://root:li13451234@localhost:5672/",
	}); err != nil {
		fmt.Printf("start sys err:%v", err)
	} else {
		fmt.Printf("start sys succ")

	}
}
