package sortslice

import (
	"fmt"
	"hash/crc32"
	"sort"
)

type uInt32Slice []uint32

func (p uInt32Slice) Len() int           { return len(p) }
func (p uInt32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p uInt32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p uInt32Slice) Sort() { sort.Sort(p) }

func GetSessionId(session []uint32) (Id uint32) {
	s := uInt32Slice(session)
	s.Sort()
	Session := ""
	for _, v := range s {
		Session = Session + fmt.Sprintf("%d", v)
	}
	Id = crc32.ChecksumIEEE([]byte(Session))
	return
}
