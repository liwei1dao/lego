package protocol

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
)

func TestMessage(t *testing.T) {
	req := NewMessage()
	req.SetVersion(1)
	req.SetMessageType(rpccore.Response)
	fmt.Printf("%b\n", req.Header[2])
	req.SetShakeHands(true)
	fmt.Printf("%b\n", req.Header[2])
	req.SetHeartbeat(true)
	fmt.Printf("%b\n", req.Header[2])
	req.SetOneway(true)
	fmt.Printf("%b\n", req.Header[2])
	req.SetCompressType(rpccore.CompressGzip)
	fmt.Printf("%b\n", req.Header[2])
	req.SetMessageStatusType(rpccore.Error)
	fmt.Printf("%b\n", req.Header[2])
	req.SetSerializeType(rpccore.ProtoBuffer)
	fmt.Printf("%b\n", req.Header[2])
	req.SetSeq(1234567890)
	req.SetFrom(&core.ServiceNode{
		Tag:  "demo",
		Type: "demo",
		Id:   "demo",
		Addr: "127.0.0.1:9852",
		Meta: map[string]string{},
	})

	m := make(map[string]string)
	m["__ID"] = "6ba7b810-9dad-11d1-80b4-00c04fd430c9"
	req.metadata = m

	payload := `{
		"A": 1,
		"B": 2,
	}
	`
	req.payload = []byte(payload)

	var buf bytes.Buffer
	_, err := req.WriteTo(&buf)
	if err != nil {
		t.Fatal(err)
	}

	res, err := Read(&buf)
	if err != nil {
		t.Fatal(err)
	}
	res.SetMessageType(rpccore.Response)

	if res.Version() != 0 {
		t.Errorf("expect 0 but got %d", res.Version())
	}

	if res.Seq() != 1234567890 {
		t.Errorf("expect 1234567890 but got %d", res.Seq())
	}

	if res.serviceMethod != "Add" || res.metadata["__ID"] != "6ba7b810-9dad-11d1-80b4-00c04fd430c9" {
		t.Errorf("got wrong metadata: %v", res.metadata)
	}

	if string(res.payload) != payload {
		t.Errorf("got wrong payload: %v", string(res.payload))
	}
}

func TestSetCompress(t *testing.T) {
	fmt.Printf("0x80 = %08b\n", 0x80)
	fmt.Printf("0x20 = %08b\n", 0x40)
	fmt.Printf("0x20 = %08b\n", 0x20)
	fmt.Printf("0x10 = %08b\n", 0x10)
	fmt.Printf("0x08 = %08b\n", 0x08)
	fmt.Printf("0x04 = %08b\n", 0x04)
	fmt.Printf("0x02 = %08b\n", 0x02)
	fmt.Printf("0x01 = %08b\n", 0x01)
	fmt.Printf("0x0C = %08b\n", 0x0E)
}
