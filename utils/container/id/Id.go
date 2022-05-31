package id

import (
	"github.com/rs/xid"
	uuid "github.com/satori/go.uuid"
)

//xid
func NewXId() string {
	id := xid.New()
	return id.String()
}

//uuid
func NewUUId() string {
	u1 := uuid.NewV4()
	return u1.String()
}
