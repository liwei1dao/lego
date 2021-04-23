package redis

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_JsonMarshal(t *testing.T) {
	result, _ := json.Marshal(100)
	fmt.Printf("结果%s \n", string(result))
}
