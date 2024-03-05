package nlp_test

import (
	"fmt"
	"testing"

	"github.com/go-ego/gse"
)

func Test_Sys(t *testing.T) {
	// 创建分词器
	var segmenter gse.Segmenter
	segmenter.LoadDict()

	// 待分析的文本
	text := "你好!"

	// 分词
	segments := segmenter.Cut(text)
	fmt.Println("分词结果：", segments)

	// 输出分词
	for _, segment := range segments {
		fmt.Printf("词: %s\n", segment)
	}

	// 获取词性
	pos := segmenter.Pos(text, false)
	fmt.Println("词性标注结果：", pos)

	// 输出分词和词性
	for i, segment := range segments {
		fmt.Printf("词: %s, 词性: %s\n", segment, pos[i])
	}

	fmt.Printf("PosStr:%s", segmenter.PosStr(pos, "你好"))
}
