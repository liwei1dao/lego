package extractor

import (
	"fmt"
	"io"
)

//提取器
type Extractor struct {
	buf              []byte
	head             int
	tail             int
	captureStartedAt int
	captured         []byte
	err              error
}

func (this *Extractor) ReadChar() (ret byte) {
	if this.head == this.tail {
		return 0
	}
	ret = this.buf[this.head]
	this.head++
	return ret
}

func (this *Extractor) NextToken() byte {
	for i := this.head; i < this.tail; i++ {
		c := this.buf[i]
		switch c {
		case ' ', '\n', '\t', '\r':
			continue
		}
		this.head = i + 1
		return c
	}
	return 0
}

func (this *Extractor) ReportError(operation string, msg string) {
	if this.err != nil {
		if this.err != io.EOF {
			return
		}
	}
	peekStart := this.head - 10
	if peekStart < 0 {
		peekStart = 0
	}
	peekEnd := this.head + 10
	if peekEnd > this.tail {
		peekEnd = this.tail
	}
	parsing := string(this.buf[peekStart:peekEnd])
	contextStart := this.head - 50
	if contextStart < 0 {
		contextStart = 0
	}
	contextEnd := this.head + 50
	if contextEnd > this.tail {
		contextEnd = this.tail
	}
	context := string(this.buf[contextStart:contextEnd])
	this.err = fmt.Errorf("%s: %s, error found in #%v byte of ...|%s|..., bigger context ...|%s|...",
		operation, msg, this.head-peekStart, parsing, context)
}
