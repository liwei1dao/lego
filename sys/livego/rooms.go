package livego

type Rooms struct {
}

func (this *Rooms) SetKey(channel string) (key string, err error) {
	return
}

func (this *Rooms) GetKey(channel string) (newKey string, err error) {
	return
}

func (this *Rooms) GetChannel(key string) (channel string, err error) {
	return
}

func (this *Rooms) DeleteChannel(channel string) (ok bool) {
	return
}

func (this *Rooms) DeleteKey(key string) (ok bool) {
	return
}
