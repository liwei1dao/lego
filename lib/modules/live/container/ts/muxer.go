package ts

const (
	tsDefaultDataLen = 184
	tsPacketLen      = 188
	h264DefaultHZ    = 90

	videoPID = 0x100
	audioPID = 0x101
	videoSID = 0xe0
	audioSID = 0xc0
)

type Muxer struct {
	videoCc  byte
	audioCc  byte
	patCc    byte
	pmtCc    byte
	pat      [tsPacketLen]byte
	pmt      [tsPacketLen]byte
	tsPacket [tsPacketLen]byte
}
