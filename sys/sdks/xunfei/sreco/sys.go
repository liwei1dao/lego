package sreco

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

func newSys(options *Options) (sys *VoiceDiscern, err error) {
	sys = &VoiceDiscern{options: options}
	return
}

type VoiceDiscern struct {
	options *Options
}

func (this VoiceDiscern) newClient(hosturl string) (conn *websocket.Conn, err error) {
	var (
		resp *http.Response
	)
	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	conn, resp, err = d.Dial(assembleAuthUrl(hosturl, this.options.ApiKey, this.options.ApiSecret), nil)
	if err != nil {
		err = fmt.Errorf("newSys:xunfei.voicediscern:%s,err:%v", readResp(resp), err)
	}
	return
}

//语音转文字
func (this VoiceDiscern) VoiceToTxt(ctx context.Context, reader io.Reader) (result string, err error) {
	var (
		st        = time.Now()
		conn      *websocket.Conn
		frameSize = 1280 //每一帧的音频大小
		buffer    = make([]byte, frameSize)
		status    = 0
		len       = 0
		msg       []byte
		decoder   Decoder
	)
	if conn, err = this.newClient("wss://iat-api.xfyun.cn/v2/iat"); err != nil {
		return
	}
	defer conn.Close()
	go func() {
		for {
			len, err = reader.Read(buffer)
			if err != nil {
				if err == io.EOF { //文件读取完了，改变status = STATUS_LAST_FRAME
					status = STATUS_LAST_FRAME
				} else {
					return
				}
			}
			select {
			case <-ctx.Done():
				err = ctx.Err()
				return
			default:
			}
			switch status {
			case STATUS_FIRST_FRAME: //发送第一帧音频，带business 参数
				frameData := map[string]interface{}{
					"common": map[string]interface{}{
						"app_id": this.options.Appid, //appid 必须带上，只需第一帧发送
					},
					"business": map[string]interface{}{ //business 参数，只需一帧发送
						"language": "zh_cn",
						"domain":   "iat",
						"accent":   "mandarin",
					},
					"data": map[string]interface{}{
						"status":   STATUS_FIRST_FRAME,
						"format":   "audio/L16;rate=16000",
						"audio":    base64.StdEncoding.EncodeToString(buffer[:len]),
						"encoding": "raw",
					},
				}
				this.options.Log.Debugf("send first")
				conn.WriteJSON(frameData)
				status = STATUS_CONTINUE_FRAME
			case STATUS_CONTINUE_FRAME:
				frameData := map[string]interface{}{
					"data": map[string]interface{}{
						"status":   STATUS_CONTINUE_FRAME,
						"format":   "audio/L16;rate=16000",
						"audio":    base64.StdEncoding.EncodeToString(buffer[:len]),
						"encoding": "raw",
					},
				}
				conn.WriteJSON(frameData)
			case STATUS_LAST_FRAME:
				frameData := map[string]interface{}{
					"data": map[string]interface{}{
						"status":   STATUS_LAST_FRAME,
						"format":   "audio/L16;rate=16000",
						"audio":    base64.StdEncoding.EncodeToString(buffer[:len]),
						"encoding": "raw",
					},
				}
				conn.WriteJSON(frameData)
				this.options.Log.Debugf("send last")
				return
			}
		}
	}()
	//获取返回的数据
	for {
		var resp = RespData{}
		_, msg, err = conn.ReadMessage()
		if err != nil {
			err = fmt.Errorf("read message error:%v", err)
			break
		}
		json.Unmarshal(msg, &resp)
		fmt.Println(resp.Data.Result.String(), resp.Sid)
		if resp.Code != 0 {
			fmt.Println(resp.Code, resp.Message, time.Since(st))
			return
		}
		decoder.Decode(&resp.Data.Result)
		if resp.Data.Status == 2 {
			fmt.Println("final:", decoder.String())
			fmt.Println(resp.Code, resp.Message, time.Since(st))
			break
		}
	}
	return
}
