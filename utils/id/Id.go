package id

import (
	"github.com/rs/xid"
)

//xid
func NewXId() string {
	id := xid.New()
	return id.String()
}
