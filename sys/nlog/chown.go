//go:build !linux
// +build !linux

package log

import (
	"os"
)

func chown(_ string, _ os.FileInfo) error {
	return nil
}
