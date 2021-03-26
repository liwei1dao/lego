package av

import "fmt"

type Info struct {
	Key   string
	URL   string
	UID   string
	Inter bool
}

func (info Info) IsInterval() bool {
	return info.Inter
}
func (info Info) String() string {
	return fmt.Sprintf("<key: %s, URL: %s, UID: %s, Inter: %v>",
		info.Key, info.URL, info.UID, info.Inter)
}

type Packet struct {
	IsAudio    bool
	IsVideo    bool
	IsMetadata bool
	TimeStamp  uint32 // dts
	StreamID   uint32
	Header     PacketHeader
	Data       []byte
}

type PacketHeader interface {
}

type GetWriter interface {
	GetWriter(Info) WriteCloser
}

type Handler interface {
	HandleReader(ReadCloser)
	HandleWriter(WriteCloser)
}

type Closer interface {
	Info() Info
	Close(error)
}

type Alive interface {
	Alive() bool
}

type CalcTime interface {
	CalcBaseTimestamp()
}

type ReadCloser interface {
	Closer
	Alive
	Read(*Packet) error
}

type WriteCloser interface {
	Closer
	Alive
	CalcTime
	Write(*Packet) error
}
