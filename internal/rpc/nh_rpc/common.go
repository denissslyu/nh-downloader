package nh_rpc

import (
	"net/http"
	"net/url"
	"nh-downloader/internal/config"
	"nh-downloader/utils/logs"
	"time"

	"github.com/denissslyu/nhentai-go"
)

var client *nhentai.Client
var transport *http.Transport

// todo config
func Init() {
	client = nhentai.NewClient()

	transport = &http.Transport{
		TLSHandshakeTimeout:   time.Second * 20,
		ExpectContinueTimeout: time.Second * 20,
		ResponseHeaderTimeout: time.Second * 20,
		IdleConnTimeout:       time.Second * 20,
	}
	client.Transport = transport

	_ = SetProxy(config.Proxy())
}

func SetProxy(proxyUrl string) error {
	logs.Info("[nh_rpc.SetProxy] setProxy:", proxyUrl)
	if proxyUrl == "" {
		transport.Proxy = nil
		client.Transport = transport
		return nil
	}

	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		logs.Error("[nh_rpc.SetProxy] parse proxyUrl failed:", err)
		return err
	}

	transport = &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}
	client.Transport = transport
	return nil
}

func SimpleSearch(filters []string, page int) (*nhentai.ComicPageData, error) {
	str := ""
	for idx, filter := range filters {
		if idx == 0 {
			str = filter
			continue
		}
		str += "+" + filter
	}
	return client.ComicByRawCondition(str, page)
}

func Info(id int) (*nhentai.ComicInfo, error) {
	info, err := client.ComicInfo(id)
	if err != nil {
		logs.Error("[nh_rpc.Info] get info failed:", err)
		return nil, err
	}
	return info, nil
}

func GetThumb(mediaId int, t string) string {
	return client.ThumbnailUrl(mediaId, t)
}

func PageUrl(mediaId, page int, t string) string {
	return client.PageUrl(mediaId, page, t)
}

func Extension(t string) string {
	return client.GetExtension(t)
}

func Get(url string) (*http.Response, error) {
	return client.Get(url)
}
