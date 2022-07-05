package kafka

import (
	"github.com/liwei1dao/lego/sys/kafka"
	"github.com/liwei1dao/lego/sys/rpcl/core"
)

func newService(sys core.ISys) (service *Service, err error) {
	service = &Service{
		sys: sys,
	}
	service.kafka, err = kafka.NewSys(
		kafka.SetVersion(sys.GetKafkaVersion()),
		kafka.SetStartType(kafka.Consumer),
		kafka.SetHosts(sys.GetKafkaAddr()),
		kafka.SetClientID(sys.GetNodePath()),
		kafka.SetTopics([]string{sys.GetNodePath()}),
		kafka.SetConsumer_Offsets_Initial(-1),
	)
	return
}

type Service struct {
	sys   core.ISys
	kafka kafka.ISys
}

func (this *Service) Start() (err error) {
	go this.run()
	return
}

func (this *Service) Close() (err error) {
	return this.kafka.Close()
}

func (this *Service) run() {
	var (
		err     error
		message *core.Message
	)
	for v := range this.kafka.Consumer_Messages() {
		message = &core.Message{}
		if err = this.sys.Decoder().Decoder(v.Value, message); err != nil {
			this.sys.Errorf("Decoder err%v", err)
			continue
		}

	}
}
