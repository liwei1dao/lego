package ssynth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	STATUS_FIRST_FRAME    = 0
	STATUS_CONTINUE_FRAME = 1
	STATUS_LAST_FRAME     = 2
)

func newSys(options *Options) (sys *VoiceSynthesis, err error) {
	sys = &VoiceSynthesis{options: options}
	return
}

type VoiceSynthesis struct {
	options *Options
}

func (this VoiceSynthesis) newClient(hosturl string) (conn *websocket.Conn, err error) {
	var (
		resp *http.Response
	)
	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	conn, resp, err = d.Dial(assembleAuthUrl(hosturl, this.options.ApiKey, this.options.ApiSecret), nil)
	if err != nil {
		err = fmt.Errorf("newSys:xunfei.voicediscern:%s,err:%v", readResp(resp), err)
	} else if resp.StatusCode != 101 {
		err = fmt.Errorf("newSys:xunfei.voicediscern:%s,err:%v", readResp(resp), err)
	}
	return
}

func (this VoiceSynthesis) TxtToVoice(ctx context.Context, srcText string, w io.Writer) (err error) {
	var (
		st         = time.Now()
		conn       *websocket.Conn
		msg        []byte
		audiobytes []byte
	)
	if conn, err = this.newClient(this.options.HostUrl); err != nil {
		return
	}
	defer conn.Close()
	frameData := map[string]interface{}{
		"common": map[string]interface{}{
			"app_id": this.options.Appid, //appid 必须带上，只需第一帧发送
		},
		"business": map[string]interface{}{ //business 参数，只需一帧发送
			"vcn": "xiaoyan",
			/*音频编码，可选值：
			raw：未压缩的pcm
			lame：mp3 (当aue=lame时需传参sfl=1)
			speex-org-wb;7： 标准开源speex（for speex_wideband，即16k）数字代表指定压缩等级（默认等级为8）
			speex-org-nb;7： 标准开源speex（for speex_narrowband，即8k）数字代表指定压缩等级（默认等级为8）
			speex;7：压缩格式，压缩等级1~10，默认为7（8k讯飞定制speex）
			speex-wb;7：压缩格式，压缩等级1~10，默认为7（16k讯飞定制speex）
			*/
			"aue":   "lame",
			"speed": 50,
			"tte":   "UTF8",
			"sfl":   1,
		},
		"data": map[string]interface{}{
			"status":   STATUS_LAST_FRAME,
			"encoding": "UTF8",
			"text":     base64.StdEncoding.EncodeToString([]byte(srcText)),
		},
	}
	fmt.Println("send first")
	conn.WriteJSON(frameData)

	//获取返回的数据

	for {
		var resp = RespData{}
		_, msg, err = conn.ReadMessage()
		if err != nil {
			fmt.Println("read message error:", err)
			break
		}
		json.Unmarshal(msg, &resp)
		//fmt.Println(string(msg))
		//fmt.Println(resp.Data.Audio, resp.Sid)
		if resp.Code != 0 {
			fmt.Println(resp.Code, resp.Message, time.Since(st))
			return
		}
		//decoder.Decode(&resp.Data.Audio)

		audiobytes, err = base64.StdEncoding.DecodeString(resp.Data.Audio)
		if err != nil {
			return
		}
		_, err = w.Write(audiobytes)
		if err != nil {
			return
		}

		if resp.Data.Status == 2 {
			fmt.Println(resp.Code, resp.Message, time.Since(st))
			break
		}
	}
	return
}
