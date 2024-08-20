package id_test

import (
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/utils/container/id"
	"github.com/sony/sonyflake"
)

func Test_IdXId(t *testing.T) {
	id := id.NewXId()
	fmt.Println(id)
}

func Test_IdUUId(t *testing.T) {
	id := id.NewUUId()
	fmt.Println(id)
}

func Test_Sonyflake(t *testing.T) {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	if sf == nil {
		fmt.Println("Sonyflake not created")
		return
	}

	id, err := sf.NextID()
	if err != nil {
		fmt.Println("Error generating ID:", err)
		return
	}

	fmt.Println("Generated ID:", id)
}
