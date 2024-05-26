package nh_rpc

import (
	"net/http"
	"nh-downloader/utils/logs"
	"time"

	"github.com/denissslyu/nhentai-go"
)

var client *nhentai.Client

// todo config
func Init() {
	client = &nhentai.Client{}

	client.Transport = &http.Transport{
		TLSHandshakeTimeout:   time.Second * 20,
		ExpectContinueTimeout: time.Second * 20,
		ResponseHeaderTimeout: time.Second * 20,
		IdleConnTimeout:       time.Second * 20,
	}
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
