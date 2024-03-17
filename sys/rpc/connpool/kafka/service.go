package kafka

import (
	"github.com/liwei1dao/lego/sys/kafka"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
)

func newService(sys rpccore.ISys, config *rpccore.Config) (service *Service, err error) {
	service = &Service{
		sys: sys,
	}
	service.kafka, err = kafka.NewSys(
		kafka.SetHosts(config.Endpoints),
		kafka.SetStartType(kafka.Consumer),
		kafka.SetClientID(sys.ServiceNode().GetNodePath()),
		kafka.SetVersion(config.Vsersion),
		kafka.SetTopics([]string{sys.ServiceNode().GetNodePath()}),
		kafka.SetConsumer_Offsets_Initial(-1),
	)
	return
}

type Service struct {
	sys   rpccore.ISys
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
	// var (
	// 	err     error
	// 	message *protocol.Message
	// )
	// for v := range this.kafka.Consumer_Messages() {
	// 	message = &protocol.Message{}
	// 	if err = this.sys.Unmarshal(v.Value, message); err != nil {
	// 		this.sys.Errorf("Decoder err%v", err)
	// 		continue
	// 	}
	// }
}
