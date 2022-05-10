package doris

import "io"

type (
	//Stream load返回消息结构体
	ResponseBody struct {
		TxnID                  int    `json:"TxnId"`
		Label                  string `json:"Label"`
		Status                 string `json:"Status"`
		Message                string `json:"Message"`
		NumberTotalRows        int    `json:"NumberTotalRows"`
		NumberLoadedRows       int    `json:"NumberLoadedRows"`
		NumberFilteredRows     int    `json:"NumberFilteredRows"`
		NumberUnselectedRows   int    `json:"NumberUnselectedRows"`
		LoadBytes              int    `json:"LoadBytes"`
		LoadTimeMs             int    `json:"LoadTimeMs"`
		BeginTxnTimeMs         int    `json:"BeginTxnTimeMs"`
		StreamLoadPutTimeMs    int    `json:"StreamLoadPutTimeMs"`
		ReadDataTimeMs         int    `json:"ReadDataTimeMs"`
		WriteDataTimeMs        int    `json:"WriteDataTimeMs"`
		CommitAndPublishTimeMs int    `json:"CommitAndPublishTimeMs"`
		ErrorURL               string `json:"ErrorURL"`
	}
	ISys interface {
		Write(tName string, body io.Reader) (err error)
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

func Write(tName string, body io.Reader) (err error) {
	return defsys.Write(tName, body)
}
