package protocol

import (
	"bytes"
	"fmt"
	"testing"

	lcore "github.com/liwei1dao/lego/sys/rpcl/core"
)

func TestMessage(t *testing.T) {
	req := NewMessage()
	req.SetVersion(1)
	fmt.Printf("%b\n", req.Header[2])
	req.SetMessageType(lcore.Response)
	fmt.Printf("%b\n", req.Header[2])
	req.SetHeartbeat(true)
	fmt.Printf("%b\n", req.Header[2])
	req.SetOneway(true)
	fmt.Printf("%b\n", req.Header[2])
	req.SetCompressType(lcore.CompressGzip)
	fmt.Printf("%b\n", req.Header[2])
	req.SetMessageStatusType(lcore.Error)
	fmt.Printf("%b\n", req.Header[2])
	req.SetSerializeType(lcore.ProtoBuffer)
	fmt.Printf("%b\n", req.Header[2])

	req.SetSeq(1234567890)

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
	res.SetMessageType(lcore.Response)

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
