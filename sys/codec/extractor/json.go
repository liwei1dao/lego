package extractor

//提取器
type Extractor struct {
	buf              []byte
	head             int
	tail             int
	captureStartedAt int
	captured         []byte
	err              error
}
