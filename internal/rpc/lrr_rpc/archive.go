package lrr_rpc

import (
	"fmt"
	"net/http"

	"nh-downloader/consts"
	"nh-downloader/utils/logs"

	"github.com/bytedance/sonic"
)

func GetMetadata(id string) (*Archive, error) {
	client := newHttpClient().
		URL(fmt.Sprintf(consts.LrrMetadataUrl, id)).
		Method(http.MethodGet)

	respBytes, err := client.send()
	if err != nil {
		logs.Error("[lrr_rpc.GetMetadata] HTTP request failed:", err)
		return nil, err
	}
	resp := &Archive{}
	err = sonic.Unmarshal(respBytes, resp)
	if err != nil {
		logs.Error("[lrr_rpc.GetMetadata] unmarshal failed:", err)
		return nil, err
	}
	return resp, nil
}

func GetThumbBytes(id string) ([]byte, error) {
	client := newHttpClient().
		URL(fmt.Sprintf(consts.LrrThumbNailUrl, id)).
		Method(http.MethodGet)

	respBytes, err := client.send()
	if err != nil {
		logs.Error("[lrr_rpc.GetThumb] HTTP request failed:", err)
		return nil, err
	}
	return respBytes, nil
}
