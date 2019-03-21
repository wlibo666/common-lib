package webutils

import (
	"testing"
)

func TestParseUrl(t *testing.T) {
	url := "www.baidu.com"
	isHttps, host, uri := ParseUrl(url)
	t.Logf("url:%s,ishttps:%v,host:%s,uri:%s", url, isHttps, host, uri)

	url = "www.baidu.com/api/v1?a=a"
	isHttps, host, uri = ParseUrl(url)
	t.Logf("url:%s,ishttps:%v,host:%s,uri:%s", url, isHttps, host, uri)

	url = "http://www.baidu.com"
	isHttps, host, uri = ParseUrl(url)
	t.Logf("url:%s,ishttps:%v,host:%s,uri:%s", url, isHttps, host, uri)

	url = "http://www.baidu.com/"
	isHttps, host, uri = ParseUrl(url)
	t.Logf("url:%s,ishttps:%v,host:%s,uri:%s", url, isHttps, host, uri)

	url = "http://www.baidu.com/?a=a"
	isHttps, host, uri = ParseUrl(url)
	t.Logf("url:%s,ishttps:%v,host:%s,uri:%s", url, isHttps, host, uri)

	url = "https://www.baidu.com"
	isHttps, host, uri = ParseUrl(url)
	t.Logf("url:%s,ishttps:%v,host:%s,uri:%s", url, isHttps, host, uri)

	url = "https://www.baidu.com/"
	isHttps, host, uri = ParseUrl(url)
	t.Logf("url:%s,ishttps:%v,host:%s,uri:%s", url, isHttps, host, uri)

	url = "https://www.baidu.com/?a=a"
	isHttps, host, uri = ParseUrl(url)
	t.Logf("url:%s,ishttps:%v,host:%s,uri:%s", url, isHttps, host, uri)

	url = "https://www.baidu.com/v1?a=a"
	isHttps, host, uri = ParseUrl(url)
	t.Logf("url:%s,ishttps:%v,host:%s,uri:%s", url, isHttps, host, uri)
}
