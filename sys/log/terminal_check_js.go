//go:build js
// +build js

package log

func isTerminal(fd int) bool {
	return false
}
