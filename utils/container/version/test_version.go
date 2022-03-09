package version

import (
	"fmt"
	"testing"
	"time"
)

func Test_Version(t *testing.T) {
	versionA := "1.2.3a "
	versionB := "1.2.3b "
	fmt.Println(CompareStrVer(versionA, versionB))
	time.LoadLocation("Local ")
	fmt.Println(time.Now())
}
