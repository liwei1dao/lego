package ftp

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func newSFtp(options Options) (sys *SFtp, err error) {
	sys = &SFtp{options: options}
	err = sys.init()
	return
}

type SFtp struct {
	options   Options
	sshClient *ssh.Client
	conn      *sftp.Client
}

func (this *SFtp) init() (err error) {
	clientConfig := &ssh.ClientConfig{
		User:            this.options.User,
		Auth:            []ssh.AuthMethod{ssh.Password(this.options.Password)},
		Timeout:         10 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if this.sshClient, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", this.options.IP, this.options.Port), clientConfig); err != nil {
		return
	}
	if this.conn, err = sftp.NewClient(this.sshClient); err != nil {
		return
	}
	return
}

func (this *SFtp) ReadDir(dir string) (files []*FileEntry, err error) {
	var (
		entries []os.FileInfo
	)
	if entries, err = this.conn.ReadDir(dir); err == nil {
		files = make([]*FileEntry, len(entries))
		for i, v := range entries {
			files[i] = &FileEntry{
				FileName: v.Name(),
				FileSize: uint64(v.Size()),
				ModTime:  v.ModTime().Unix(),
				IsDir:    v.IsDir(),
			}
		}
	}
	return
}

func (this *SFtp) MakeDir(path string) error {
	return this.conn.MkdirAll(path)
}

func (this *SFtp) RemoveDir(path string) error {
	return this.conn.RemoveDirectory(path)
}

func (this *SFtp) Close() (err error) {
	this.sshClient.Close()
	return
}
