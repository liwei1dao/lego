package id_test

import (
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/utils/container/id"
)

func Test_IdXId(t *testing.T) {
	id := id.NewXId()
	fmt.Println(id)
}

func Test_IdUUId(t *testing.T) {
	id := id.NewUUId()
	fmt.Println(id)
}
