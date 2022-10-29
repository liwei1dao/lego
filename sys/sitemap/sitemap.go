package sitemap

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
	"github.com/liwei1dao/lego/utils/codec"
)

var (
	TooMuchLinksError = errors.New("单个sitemap文件过多")
)

type UrlSet struct {
	*base
	XMLName    xml.Name `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
	XMLNSVideo string   `xml:"xmlns:video,attr,omitempty"`
	XMLNSImage string   `xml:"xmlns:image,attr,omitempty"`
	XMLNSNews  string   `xml:"xmlns:news,attr,omitempty"`
	Token      []xml.Token
}

func newSys(options *Options) (sys *Sitemap, err error) {
	sys = &Sitemap{
		urlSet: &UrlSet{
			base: &base{},
		},
		urls:    map[string]int{},
		options: options,
	}
	err = sys.load()
	for i, v := range sys.urlSet.Token {
		url := v.(*Url)
		sys.urls[url.Loc] = i
	}
	return
}

type Sitemap struct {
	urlSet  *UrlSet
	urls    map[string]int
	options *Options
}

func (this *Sitemap) AppendUrl(url *Url) {
	if !strings.HasPrefix(url.Loc, "http") {
		url.Loc = strings.TrimRight(this.options.DefaultHost, "/") + strings.TrimLeft(url.Loc, "/")
	}
	this.urlSet.setNs(url.xmlns)
	if i, ok := this.urls[url.Loc]; ok {
		this.urlSet.Token[i] = url
	} else {
		this.urlSet.Token = append(this.urlSet.Token, url)
		this.urls[url.Loc] = len(this.urlSet.Token) - 1
	}
}

func (this *Sitemap) GetUrls() (urls []*Url) {
	urls = make([]*Url, 0)
	for _, v := range this.urlSet.Token {
		url := v.(*Url)
		urls = append(urls, url)
	}
	return
}

func (this *Sitemap) SetUrls(urls []*Url) {
	this.urlSet.Token = make([]xml.Token, 0, len(urls))
	for _, v := range urls {
		this.urlSet.Token = append(this.urlSet.Token, v)
	}
}

func (this *Sitemap) ToXml() ([]byte, error) {
	if ImageXmlNS&this.urlSet.xmlns == ImageXmlNS {
		this.urlSet.XMLNSImage = "http://www.google.com/schemas/sitemap-image/1.1"
	}
	if VideoXmlNS&this.urlSet.xmlns == VideoXmlNS {
		this.urlSet.XMLNSVideo = "http://www.google.com/schemas/sitemap-video/1.1"
	}
	if NewsXmlNS&this.urlSet.xmlns == NewsXmlNS {
		this.urlSet.XMLNSNews = "http://www.google.com/schemas/sitemap-news/0.9"
	}
	if len(this.urlSet.Token) > this.options.MaxLinks {
		return nil, TooMuchLinksError
	}
	var (
		data []byte
		err  error
		buf  bytes.Buffer
	)
	if this.options.Pretty {
		buf.Write([]byte(xml.Header))
		data, err = xml.MarshalIndent(this.urlSet, "", "  ")
	} else {
		buf.Write([]byte(strings.Trim(xml.Header, "\n")))
		data, err = xml.Marshal(this.urlSet)
	}
	if err != nil {
		return nil, err
	}
	buf.Write(data)
	return buf.Bytes(), nil
}

// filename 生成sitemap文件名
func (this *Sitemap) Storage() (filename string, err error) {
	var (
		data []byte
	)
	data, err = this.ToXml()
	if err != nil {
		return
	}
	if err = os.MkdirAll(filepath.Dir(this.options.Filename), 0755); err != nil {
		return
	}
	if this.options.Compress {
		basename := strings.TrimRight(this.options.Filename, path.Ext(this.options.Filename))
		filename = basename + ".xml.gz"
		if fd, err := os.OpenFile(this.options.Filename, os.O_WRONLY|os.O_CREATE, 0666); err != nil {
			return "", err
		} else {
			gw := gzip.NewWriter(fd)
			if _, err := gw.Write(data); err != nil {
				return "", err
			}
			_ = fd.Close()
			_ = gw.Close()
			return filename, nil
		}
	} else {
		return this.options.Filename, ioutil.WriteFile(this.options.Filename, data, 0666)
	}
}

func (this *Sitemap) load() (err error) {
	if _, err := os.Stat(this.options.Filename); err != nil {
		return nil
	}
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(this.options.Filename); err != nil {
		panic(err)
	}

	root := doc.SelectElement("urlset")
	for _, url := range root.SelectElements("url") {
		_url := NewUrl()
		for _, v := range url.ChildElements() {
			fmt.Println("url :", v.Tag)
			switch v.Tag {
			case "loc":
				_url.Loc = v.Text()
			case "lastmod":
				_url.LastMod = v.Text()
			case "changefreq":
				_url.ChangeFreq = ChangeFreq(v.Text())
			case "priority":
				_url.Priority = codec.StringToFloat64(v.Text())
			case "Image":
				_url.AppendImage(NewImageForXml(v))
			case "Video":
				_url.AppendVideo(NewVideoForXml(v))
			case "News":
				_url.AppendNews(NewNewsForXml(v))
			}
		}
		this.AppendUrl(_url)
	}
	return
}

type Map struct {
	XMLName xml.Name `xml:"sitemap"`
	Loc     string   `xml:"loc"`
}

type siteMapIndex struct {
	XMLName xml.Name `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 sitemapindex"`
	SiteMap []Map
}

func NewSiteMapIndex() *siteMapIndex {
	return &siteMapIndex{
		SiteMap: make([]Map, 0),
	}
}

func (this *siteMapIndex) Append(loc string) {
	m := Map{
		Loc: loc,
	}
	this.SiteMap = append(this.SiteMap, m)
}

func (this *siteMapIndex) ToXml() ([]byte, error) {
	var (
		data []byte
		err  error
		buf  bytes.Buffer
	)
	buf.Write([]byte(xml.Header))
	data, err = xml.MarshalIndent(this, "", "  ")
	if err != nil {
		return nil, err
	}
	buf.Write(data)
	return buf.Bytes(), nil
}

func (this *siteMapIndex) Storage(filepath string) (filename string, err error) {
	if path.Ext(filepath) != ".xml" {
		return "", errors.New("建议以.xml作为文件扩展名")
	}
	var data []byte
	if data, err = this.ToXml(); err != nil {
		return
	}
	err = ioutil.WriteFile(filepath, data, 0666)
	filename = path.Base(filepath)
	return
}
