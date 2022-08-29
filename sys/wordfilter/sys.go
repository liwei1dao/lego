package wordfilter

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

func newSys(options *Options) (sys *Sys, err error) {
	sys = &Sys{
		options: options,
		trie:    NewTrie(),
		noise:   regexp.MustCompile(`[\|\s&%$@*]+`),
	}
	return
}

type Sys struct {
	options *Options
	trie    *Trie
	noise   *regexp.Regexp
}

// UpdateNoisePattern 更新去噪模式
func (Sys *Sys) UpdateNoisePattern(pattern string) {
	Sys.noise = regexp.MustCompile(pattern)
}

// LoadWordDict 加载敏感词字典
func (Sys *Sys) LoadWordDict(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return Sys.Load(f)
}

// LoadNetWordDict 加载网络敏感词字典
func (Sys *Sys) LoadNetWordDict(url string) error {
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	rsp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	return Sys.Load(rsp.Body)
}

// Load common method to add words
func (Sys *Sys) Load(rd io.Reader) error {
	buf := bufio.NewReader(rd)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		Sys.trie.Add(string(line))
	}
	return nil
}

// AddWord 添加敏感词
func (Sys *Sys) AddWord(words ...string) {
	Sys.trie.Add(words...)
}

// DelWord 删除敏感词
func (Sys *Sys) DelWord(words ...string) {
	Sys.trie.Del(words...)
}

// Sys 过滤敏感词
func (Sys *Sys) Filter(text string) string {
	return Sys.trie.Filter(text)
}

// Replace 和谐敏感词
func (Sys *Sys) Replace(text string, repl rune) string {
	return Sys.trie.Replace(text, repl)
}

// FindIn 检测敏感词
func (Sys *Sys) FindIn(text string) (bool, string) {
	text = Sys.RemoveNoise(text)
	return Sys.trie.FindIn(text)
}

// FindAll 找到所有匹配词
func (Sys *Sys) FindAll(text string) []string {
	return Sys.trie.FindAll(text)
}

// Validate 检测字符串是否合法
func (Sys *Sys) Validate(text string) (bool, string) {
	text = Sys.RemoveNoise(text)
	return Sys.trie.Validate(text)
}

// RemoveNoise 去除空格等噪音
func (Sys *Sys) RemoveNoise(text string) string {
	return Sys.noise.ReplaceAllString(text, "")
}
