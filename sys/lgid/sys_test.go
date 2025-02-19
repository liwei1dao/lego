package lgid_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/bwmarrin/snowflake"
)

func Test_sys(t *testing.T) {
	// 初始化雪花节点，假设节点 ID 为 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatalf("error initializing snowflake node: %v", err)
	}

	// 生成唯一 ID
	id := node.Generate()
	fmt.Printf("Unique ID: %d\n", id.Int64())
}
