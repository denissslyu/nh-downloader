package lrr_serv

import (
	"slices"
	"strings"

	"nh-downloader/consts"
	"nh-downloader/internal/model"
	"nh-downloader/internal/rpc/lrr_rpc"
	"nh-downloader/utils"
	"nh-downloader/utils/logs"

	"github.com/spf13/cast"
)

func GetItem(id string) (*model.Item, error) {
	arch, err := lrr_rpc.GetMetadata(id)
	if err != nil {
		return nil, err
	}
	if arch == nil {
		return nil, nil
	}
	item := convLrrArchiveToItem(arch)
	return item, nil
}

// SimpleSearchItems search items with combined filter words
func SimpleSearchItems(option *model.SimpleSearchOption) (*model.SimpleSearchResult, error) {
	req := convSimpleSearchOptionToSearchReq(option)
	resp, err := lrr_rpc.Search(req)
	if err != nil {
		return nil, err
	}
	if resp == nil || len(resp.Data) == 0 {
		return nil, nil
	}
	result := &model.SimpleSearchResult{
		Option:    option,
		Items:     make([]*model.Item, 0, len(resp.Data)),
		TotolPage: resp.RecordsFiltered/consts.LrrPageSize + 1,
	}

	for _, arch := range resp.Data {
		result.Items = append(result.Items, convLrrArchiveToItem(arch))
	}
	return result, nil
}

func convSimpleSearchOptionToSearchReq(option *model.SimpleSearchOption) *lrr_rpc.SearchReq {
	req := &lrr_rpc.SearchReq{
		SortBy: option.Sort,
		Order:  option.Order,
	}
	str := ""
	for idx, filter := range option.Filters {
		if idx == 0 {
			str = filter
			continue
		}
		str += "," + filter
	}
	req.Filter = str
	if option.Page > 1 {
		req.Start = cast.ToString((option.Page - 1) * consts.LrrPageSize)
	}
	return req
}

func convLrrArchiveToItem(arch *lrr_rpc.Archive) *model.Item {
	item := &model.Item{
		Title: model.Title{
			Main:       arch.Title,
			MainPretty: utils.GetPrettyTitle(arch.Title),
		},
		Id:     cast.ToString(arch.Arcid),
		Source: consts.Lanraragi,
		Pages:  arch.Pagecount,
	}
	tags := strings.Split(arch.Tags, ",")
	for _, tag := range tags {
		parts := strings.SplitN(tag, ":", 2)
		if len(parts) < 2 {
			item.Tags = append(item.Tags, model.Tag{Type: consts.TagType_Tag, Value: strings.TrimSpace(tag)})
			continue
		}
		tagType := strings.ToLower(strings.TrimSpace(parts[0]))
		tagValue := strings.TrimSpace(parts[1])
		if tagType == consts.TagType_DateAdded {
			item.UploadDate = cast.ToInt(tagValue)
		}
		if tagType == consts.TagType_Pages {
			continue
		}
		if newTagType, ok := consts.TagTypeReplaceMap[tagType]; ok {
			tagType = newTagType
		}
		if !slices.Contains(consts.TagTypes, tagType) {
			logs.Debug("strange tag type:", tagType)
		}
		item.Tags = append(item.Tags, model.Tag{Type: tagType, Value: tagValue})

	}
	item.SetLangByTag()

	return item
}
