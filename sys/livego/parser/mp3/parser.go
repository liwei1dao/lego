package mp3

import "fmt"

// sampling_frequency - indicates the sampling frequency, according to the following table.
// '00' 44.1 kHz
// '01' 48 kHz
// '10' 32 kHz
// '11' reserved
var mp3Rates = []int{44100, 48000, 32000}
var (
	errMp3DataInvalid = fmt.Errorf("mp3data  invalid")
	errIndexInvalid   = fmt.Errorf("invalid rate index")
)

func NewParser() *Parser {
	return &Parser{}
}

type Parser struct {
	samplingFrequency int
}

func (this *Parser) SampleRate() int {
	if this.samplingFrequency == 0 {
		this.samplingFrequency = 44100
	}
	return this.samplingFrequency
}

func (this *Parser) Parse(src []byte) error {
	if len(src) < 3 {
		return errMp3DataInvalid
	}
	index := (src[2] >> 2) & 0x3
	if index <= byte(len(mp3Rates)-1) {
		this.samplingFrequency = mp3Rates[index]
		return nil
	}
	return errIndexInvalid
}
