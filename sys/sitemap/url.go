package sitemap

import (
	"encoding/xml"
	"time"
)

type ChangeFreq string
type InvalidPriorityError struct {
	msg string
}

func (e *InvalidPriorityError) Error() string {
	return e.msg
}

const (
	Always  ChangeFreq = "always"
	Hourly  ChangeFreq = "hourly"
	Daily   ChangeFreq = "daily"
	Weekly  ChangeFreq = "weekly"
	Monthly ChangeFreq = "monthly"
	Yearly  ChangeFreq = "yearly"
	Never   ChangeFreq = "never"
)

func NewUrl() *Url {
	return &Url{
		base:       &base{},
		Loc:        "",
		LastMod:    "",
		ChangeFreq: "",
		Priority:   0,
	}
}

type Url struct {
	*base
	XMLName    xml.Name   `xml:"url"`
	Loc        string     `xml:"loc"`                  //链接地址
	LastMod    string     `xml:"lastmod,omitempty"`    //最后一次更新时间
	ChangeFreq ChangeFreq `xml:"changefreq,omitempty"` //变化周期
	Priority   float64    `xml:"priority,omitempty"`   //优先级
	Token      []xml.Token
}

// 网址
func (u *Url) SetLoc(loc string) *Url {
	u.Loc = loc
	return u
}

// 最后一次修改时间
func (u *Url) SetLastmod(lastMod time.Time) *Url {
	u.LastMod = lastMod.Format(time.RFC3339)
	return u
}

// 更新频率
func (u *Url) SetChangefreq(freq ChangeFreq) *Url {
	u.ChangeFreq = freq
	return u
}

// 网页优先级
func (u *Url) SetPriority(priority float64) *Url {
	if priority < 0 || priority > 1 {
		panic(InvalidPriorityError{"Valid values range from 0.0 to 1.0"})
	}
	u.Priority = priority
	return u
}

// 对于单个网页上的多个视频，为该网页创建一个 <loc> 标记，并为该网页上的每个视频创建一个子级 <video> 元素。
func (u *Url) AppendVideo(video *Video) {
	u.setNs(VideoXmlNS)
	u.Token = append(u.Token, video)
}

// 对于单个网页上的多个图片，每个 <url> 标记最多可包含 1000 个 <image:image> 标记。
func (u *Url) AppendImage(image *Image) {
	u.setNs(ImageXmlNS)
	u.Token = append(u.Token, image)
}

func (u *Url) AppendNews(news *News) {
	u.setNs(NewsXmlNS)
	u.Token = append(u.Token, news)
}
