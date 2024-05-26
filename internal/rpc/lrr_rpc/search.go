package lrr_rpc

import (
	"net/http"

	"nh-downloader/consts"
	"nh-downloader/utils"
	"nh-downloader/utils/logs"

	"github.com/bytedance/sonic"
)

type SearchReq struct {
	Category string
	Filter   string
	Start    string
	SortBy   string
	Order    string
}

type SearchResp struct {
	Data            []*Archive `json:"data"`
	Draw            int        `json:"draw"`
	RecordsFiltered int        `json:"recordsFiltered"`
	RecordsTotal    int        `json:"recordsTotal"`
}

func Search(req *SearchReq) (*SearchResp, error) {
	client := newHttpClient().
		URL(consts.LrrSearchUrl).
		Method(http.MethodGet).
		Params(utils.StructToMap(req))

	respBytes, err := client.send()
	if err != nil {
		logs.Error("[lrr_rpc.Search] HTTP request failed:", err)
		return nil, err
	}
	resp := &SearchResp{}
	err = sonic.Unmarshal(respBytes, resp)
	if err != nil {
		logs.Error("[lrr_rpc.Search] unmarshal failed:", err)
		return nil, err
	}
	return resp, nil
}
