package flietools

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/liwei1dao/lego/utils/codec/json"
)

//判断文件或文件夹是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

//获取去除文件后缀的文件名称
func GetFileNameSubSuffix(filepath string) string {
	var fileSuffix string
	fileSuffix = path.Ext(filepath)
	var filenameOnly string
	filenameOnly = strings.TrimSuffix(filepath, fileSuffix)
	return filenameOnly
}

//读取json文件到结构体中 参数必须是指针
func ReadJsonFileToStruct(path string, d interface{}) error {
	var data []byte
	buf := new(bytes.Buffer)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadSlice('\n')
		if err != nil {
			if len(line) > 0 {
				buf.Write(line)
			}
			break
		}
		if !strings.HasPrefix(strings.TrimLeft(string(line), "\t "), "//") {
			buf.Write(line)
		}
	}
	data = buf.Bytes()
	return json.Unmarshal(data, d)
}

//将数据写入json文件中
func WrietStructToJsonFile(path string, d interface{}) error {
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0644)
	return err
}

//创建目录文件
func CreateDirectory(logpath string) error {
	logdir := string(logpath[0:strings.LastIndex(logpath, "/")])
	if !IsExist(logdir) {
		err := os.MkdirAll(logdir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("创建目录文件失败 1" + err.Error())
		}
	}
	return nil
}
