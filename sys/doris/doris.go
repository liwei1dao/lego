package doris

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/liwei1dao/lego/utils/codec/json"

	"github.com/liwei1dao/lego/utils/container/id"
)

func newSys(options Options) (sys *Doris, err error) {
	sys = &Doris{options: options}
	return
}

type Doris struct {
	options Options
}

//发送数据
func (this *Doris) Write(tName string, body io.Reader) (err error) {
	client := &http.Client{}
	reqest, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s:%d/api/%s/%s/_stream_load", this.options.IP, this.options.Port, this.options.DBname, tName), body)
	//增加header选项
	reqest.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(this.options.User+":"+this.options.Password)))
	reqest.Header.Add("EXPECT", "100-continue")
	reqest.Header.Add("label", id.NewXId())
	reqest.Header.Add("column_separator", this.options.Separator)
	//处理返回结果
	response, err := client.Do(reqest)
	if err == nil {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		responseBody := ResponseBody{}
		if err = json.Unmarshal(body, &responseBody); err != nil {
			return
		}
		if responseBody.Status != "Success" {
			err = fmt.Errorf("responseBody Status is:%s", responseBody.Status)
		}
	}
	return
}
