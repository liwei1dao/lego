package ip

import (
	"io/ioutil"
	"net"
	"net/http"

	"github.com/liwei1dao/lego/utils/codec/json"

	"github.com/axgle/mahonia"
)

type IPInfo struct {
	IP          string `json:"ip"`
	Pro         string `json:"pro"`
	ProCode     string `json:"proCode"`
	City        string `json:"city"`
	CityCode    string `json:"cityCode"`
	Region      string `json:"region"`
	RegionCode  string `json:"regionCode"`
	Addr        string `json:"addr"`
	RegionNames string `json:"regionNames"`
	Err         string `json:"err"`
}

//获取以太网IP
func GetEthernetInfo() *IPInfo {
	url := "http://whois.pconline.com.cn/ipJson.jsp?json=true"
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	bodystr := mahonia.NewDecoder("gbk").ConvertString(string(body))
	var result IPInfo
	if err := json.Unmarshal([]byte(bodystr), &result); err != nil {
		return nil
	}
	return &result
}

//获取本地Ip
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
