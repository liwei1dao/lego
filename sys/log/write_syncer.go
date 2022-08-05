package log

import (
	"io"
	"sync"

	"go.uber.org/multierr"
)

func AddWrites(w ...IWrite) IWrite {
	ws := make(writeTree, 0, len(w))
	for _, v := range w {
		ws = append(ws, v)
	}
	return ws
}

type writeTree []IWrite

func (mc writeTree) WriteTo(bs []byte) error {
	var err error
	for i := range mc {
		err = multierr.Append(err, mc[i].WriteTo(bs))
	}
	return err
}

func (mc writeTree) Sync() error {
	var err error
	for i := range mc {
		err = multierr.Append(err, mc[i].Sync())
	}
	return err
}

type IWriteSyncer interface {
	io.Writer
	Sync() error
}
type IWrite interface {
	WriteTo(bs []byte) error
	Sync() error
}

func AddSync(w io.Writer) IWrite {
	switch w := w.(type) {
	case IWrite:
		return w
	default:
		return writerWrapper{w}
	}
}

type writerWrapper struct {
	io.Writer
}

func (w writerWrapper) WriteTo(bs []byte) error {
	_, err := w.Write(bs)
	return err
}
func (w writerWrapper) Sync() error {
	return nil
}

func Lock(w IWriteSyncer) IWrite {
	return &lockedWriteSyncer{w: w}
}

type lockedWriteSyncer struct {
	sync.Mutex
	w IWriteSyncer
}

func (s *lockedWriteSyncer) WriteTo(bs []byte) error {
	s.Lock()
	_, err := s.w.Write(bs)
	s.Unlock()
	return err
}

func (s *lockedWriteSyncer) Sync() error {
	s.Lock()
	err := s.w.Sync()
	s.Unlock()
	return err
}
