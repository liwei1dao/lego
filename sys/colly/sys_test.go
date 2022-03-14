package colly_test

import (
	"fmt"
	"testing"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
)

func Test_sys(t *testing.T) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36"),
		colly.AllowURLRevisit(),
	)

	if p, err := proxy.RoundRobinProxySwitcher(
		"socks5://127.0.0.1:1080",
		"http://127.0.0.1:1087",
	); err == nil {
		c.SetProxyFunc(p)
	}
	c.OnRequest(func(r *colly.Request) {
		// r.Headers.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		// r.Headers.Set("accept-encoding", "gzip, deflate, br")
		// r.Headers.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
		// r.Headers.Set("cache-control", "max-age=0")
		// r.Headers.Set("cookie", "PHPSESSID=b00jn67sruelb4npphlvigkhk7; history=%5B%7B%22name%22%3A%22%E5%A4%A9%E7%9B%9B%E9%95%BF%E6%AD%8C%22%2C%22pic%22%3A%22https%3A%2F%2Fae01.alicdn.com%2Fkf%2FU4146e00599ec4bc0808ce8f9759a4874v.jpg%22%2C%22link%22%3A%22%2Fvod%2Fplay%2Fid%2F9201%2Fsid%2F2%2Fnid%2F4.html%22%2C%22part%22%3A%22%E7%AC%AC4%E9%9B%86%22%7D%5D; Hm_lvt_c96a948c88ef7639ead13a51733f3057=1629864853,1630994953; Hm_lpvt_c96a948c88ef7639ead13a51733f3057=1630994953")
		// r.Headers.Set("if-modified-since", "Tue, 07 Sep 2021 08:31:04 GMT")
		// r.Headers.Set("sec-ch-ua", "'Chromium';v='92', ' Not A;Brand';v='99', 'Google Chrome';v='92'")
		// r.Headers.Set("sec-ch-ua-mobile", "?0")
		// r.Headers.Set("sec-fetch-dest", "document")
		// r.Headers.Set("sec-fetch-mode", "navigate")
		// r.Headers.Set("sec-fetch-site", "none")
		// r.Headers.Set("sec-fetch-user", "?1")
		// r.Headers.Set("upgrade-insecure-requests", "1")
		// r.Headers.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")
		fmt.Printf("OnRequest %v\n", r)
	})
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Printf("OnError %v\n", err)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("OnResponse %v\n", string(r.Body))
	})
	c.OnScraped(func(r *colly.Response) {
		fmt.Printf("OnScraped %v\n", r)
	})
	c.Visit("https://www.tangrenjie.tv/api.php/provide/vod/?ac=list&h=24")
}

// func Test_AppleSMC(t *testing.T) {
// 	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct)
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
// 		os.Exit(1)
// 	}
// 	// setup a http client
// 	httpTransport := &http.Transport{}
// 	httpClient := &http.Client{Transport: httpTransport}
// 	// set our socks5 as the dialer
// 	httpTransport.Dial = dialer.Dial
// 	if resp, err := httpClient.Get("https://www.tangrenjie.tv/api.php/provide/vod/?ac=list&h=24"); err != nil {
// 		fmt.Printf("err:%v", err)
// 	} else {
// 		defer resp.Body.Close()
// 		body, _ := ioutil.ReadAll(resp.Body)
// 		fmt.Printf("%s\n", body)
// 	}
// }
