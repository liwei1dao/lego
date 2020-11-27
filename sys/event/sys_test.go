package event_test

import (
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/core"
)

func Test_sys(t *testing.T) {
	if err := OnInit(nil); err == nil {
		Register(core.Event_Key("TestEvent"), func() {
			fmt.Printf("TestEvent TriggerEvent")
		})
		TriggerEvent(core.Event_Key("TestEvent"))
	}
}
