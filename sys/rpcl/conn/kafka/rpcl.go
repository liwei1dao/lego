package kafka

type RPCL struct {
	service *Service
	clients map[string]*Client
}
