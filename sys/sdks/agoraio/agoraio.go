package agoraio

import (
	"time"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/RtcTokenBuilder"
)

func newSys(options Options) (sys *Agoraio, err error) {
	sys = &Agoraio{options: options}
	return
}

type Agoraio struct {
	options Options
}

func (this *Agoraio) CreateToken(uid uint32, channelName string) (token string, err error) {
	expireTimestamp := uint32(time.Now().UTC().Unix()) + this.options.ExpireTimeInSeconds
	token, err = rtctokenbuilder.BuildTokenWithUID(this.options.AppID, this.options.AppCertificate, channelName, uid, rtctokenbuilder.RoleAttendee, expireTimestamp)
	return
}
