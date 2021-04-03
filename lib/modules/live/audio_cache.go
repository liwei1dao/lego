package live

import "bytes"

const (
	cache_max_frames byte = 6
	audio_cache_len  int  = 10 * 1024
)

type audioCache struct {
	soundFormat byte
	num         byte
	offset      int
	pts         uint64
	buf         *bytes.Buffer
}
