package render

import (
	"net/http"
)

type Render interface {
	// Render 使用自定义 ContentType 写入数据。
	Render(http.ResponseWriter) error
	// WriteContentType 写入自定义 ContentType。
	WriteContentType(w http.ResponseWriter)
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
