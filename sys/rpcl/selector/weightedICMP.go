package selector

import (
	"context"

	"github.com/liwei1dao/lego/sys/rpcl/core"
)

//权重加轮询选择器
type weightedICMPSelector struct {
}

func newWeightedICMPSelector(servers map[string]string) core.ISelector {
	return &weightedICMPSelector{}
}

func (s weightedICMPSelector) Select(ctx context.Context, route core.IRoute, serviceMethod string) string {

	return ""
}

func (s *weightedICMPSelector) UpdateServer(servers map[string]string) {

}

// func createICMPWeighted(servers map[string]string) []*Weighted {
// 	var ss = make([]*Weighted, 0, len(servers))
// 	for k := range servers {
// 		w := &Weighted{Server: k, Weight: 1, EffectiveWeight: 1}
// 		server := strings.Split(k, "@")
// 		host, _, _ := net.SplitHostPort(server[1])
// 		rtt, _ := Ping(host)
// 		rtt = CalculateWeight(rtt)
// 		w.Weight = rtt
// 		w.EffectiveWeight = rtt
// 		ss = append(ss, w)
// 	}
// 	return ss
// }

// // Ping gets network traffic by ICMP
// func Ping(host string) (rtt int, err error) {
// 	rtt = 1000 //default and timeout is 1000 ms

// 	pinger, err := ping.NewPinger(host)
// 	if err != nil {
// 		return rtt, err
// 	}
// 	pinger.Count = 3
// 	pinger.Timeout = 3 * time.Second
// 	err = pinger.Run()
// 	if err != nil {
// 		return rtt, err
// 	}
// 	stats := pinger.Statistics()
// 	// ping failed
// 	if len(stats.Rtts) == 0 {
// 		return rtt, err
// 	}
// 	rtt = int(stats.AvgRtt) / 1e6

// 	return rtt, err
// }

// func CalculateWeight(rtt int) int {
// 	switch {
// 	case rtt >= 0 && rtt <= 10:
// 		return 191
// 	case rtt > 10 && rtt <= 200:
// 		return 201 - rtt
// 	case rtt > 100 && rtt < 1000:
// 		return 1
// 	default:
// 		return 0
// 	}
// }
