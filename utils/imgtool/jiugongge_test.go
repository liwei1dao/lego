package imgtool

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/disintegration/imaging"
)

func TestBuildJiuGong(t *testing.T) {
	leng := 5
	src := make([]io.Reader, leng)
	f, _ := os.Create("./text.jpg")
	defer f.Close()
	for i := 0; i < leng; i++ {
		if file, err := os.Open(fmt.Sprintf("./head00%d.jpg", i+1)); err == nil {
			src[i] = file
		}
	}
	BuildJiuGong(src, f, imaging.JPEG)
}
