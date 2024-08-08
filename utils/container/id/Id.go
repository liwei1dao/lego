package id

import (
	"fmt"
	"time"

	"github.com/rs/xid"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/exp/rand"
)

// xid
func NewXId() string {
	id := xid.New()
	return id.String()
}

// uuid
func NewUUId() string {
	u1 := uuid.NewV4()
	return u1.String()
}

func NewNanoid(){
	
}