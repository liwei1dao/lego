package core

import "bufio"

type ReadWriter struct {
	*bufio.ReadWriter
	readError  error
	writeError error
}
