package sitemap_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/liwei1dao/lego/sys/sitemap"
)

//系统测试
func Test_Sys(t *testing.T) {
	if sys, err := sitemap.NewSys(
		sitemap.SetPretty(true),
		sitemap.SetFilename("./sitemap.xml"),
	); err != nil {
		t.Log(err)
	} else {
		url := sitemap.NewUrl()
		url.SetLoc("https://www.douyacun.com/")
		url.SetLastmod(time.Now())
		url.SetChangefreq(sitemap.Daily)
		url.SetPriority(float64(3.234) / float64(10))
		sys.AppendUrl(url)
		bt, err := sys.ToXml()
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		fmt.Printf("%s", bt)
		sys.Storage()
	}
}

func TestNewVideo(t *testing.T) {
	if sys, err := sitemap.NewSys(
		sitemap.SetPretty(true),
		sitemap.SetFilename("./sitemap.xml"),
	); err != nil {
		t.Log(err)
	} else {
		url := sitemap.NewUrl()
		url.Loc = "https://www.douyacun.com"

		video := sitemap.NewVideo().
			SetThumbnailLoc("http://www.example.com/thumbs/123.jpg").
			SetTitle("适合夏季的烧烤排餐").
			SetDescription("小安教您如何每次都能烤出美味牛排").
			SetContentLoc("http://streamserver.example.com/video123.mp4").
			SetPlayerLoc("http://www.example.com/videoplayer.php?video=123", true).
			SetDuration(600*time.Second).
			SetExpirationDate(time.Now().Add(time.Hour*24)).
			SetRating(4.2).
			SetViewCount(12345).
			SetPublicationDate(time.Now().Add(-time.Hour*24)).
			SetFamilyFriendly(true).
			SetRestriction([]string{"IE", "GB", "US", "CN", "CA"}, true).
			SetPrice(6.99, "EUR", true, true).
			SetRequiresSubscription(true).
			SetUploader("GrillyMcGrillerson", "http://www.example.com/users/grillymcgrillerson").
			SetLive(true)

		url.AppendVideo(video)
		sys.AppendUrl(url)
		bt, err := sys.ToXml()
		if err != nil {
			log.Printf("%v", err)
			return
		}
		fmt.Printf("%s", bt)
		sys.Storage()
	}

}

func TestNewImage(t *testing.T) {
	if sys, err := sitemap.NewSys(
		sitemap.SetPretty(true),
		sitemap.SetFilename("./sitemap.xml"),
	); err != nil {
		t.Log(err)
	} else {
		url := sitemap.NewUrl()
		url.SetLoc("https://www.douyacun.com/show_image")
		// 注意这里url.SetLoc设置网页的网址
		// image.SetLoc 设置的是图片访问路径，image的域名和网址域名不一致也可以
		_image := sitemap.NewImage().
			SetLoc("https://www.douyacun.com/image1.jpg").
			SetTitle("example").
			SetCaption("图片的说明。").
			SetLicense("https://www.douyacun.com")
		url.AppendImage(_image)

		_image = sitemap.NewImage().
			SetLoc("https://www.douyacun.com/image2.jpg").
			SetTitle("example").
			SetCaption("图片的说明。").
			SetLicense("https://www.douyacun.com")
		url.AppendImage(_image)

		sys.AppendUrl(url)
		data, err := sys.ToXml()
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		fmt.Printf("%s", data)
	}

}

func TestNewNews(t *testing.T) {
	if sys, err := sitemap.NewSys(
		sitemap.SetPretty(true),
		sitemap.SetFilename("./sitemap.xml"),
	); err != nil {
		t.Log(err)
	} else {
		url := sitemap.NewUrl()
		url.SetLoc("https://www.douyacun.com/business/article55.html")
		url.AppendNews(sitemap.NewNews().SetName("《示例时报》").
			SetTitle("公司 A 和 B 正在进行合并谈判").
			SetLanguage("zh-cn").
			SetPublicationDate(time.Now()))
		sys.AppendUrl(url)
		data, err := sys.ToXml()
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		fmt.Printf("%s", data)
	}
}
