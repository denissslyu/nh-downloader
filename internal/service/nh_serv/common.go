package nh_serv

import (
	"slices"
	"strings"

	"nh-downloader/consts"
	"nh-downloader/internal/model"
	"nh-downloader/internal/rpc/nh_rpc"
	"nh-downloader/utils"
	"nh-downloader/utils/logs"

	"github.com/denissslyu/nhentai-go"
	"github.com/spf13/cast"
)

func GetItem(idStr string) (*model.Item, error) {
	id, err := cast.ToIntE(idStr)
	if err != nil {
		logs.Error("invalid id:", idStr)
		return nil, err
	}
	info, err := nh_rpc.Info(id)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}
	item := convNhInfoToItem(info)
	return item, nil
}

// SimpleSearchItems search items with combined filter words
func SimpleSearchItems(option *model.SimpleSearchOption) (*model.SimpleSearchResult, error) {
	data, err := nh_rpc.SimpleSearch(option.Filters, option.Page)
	if err != nil {
		return nil, err
	}
	if data == nil || len(data.Records) == 0 {
		return nil, nil
	}
	result := &model.SimpleSearchResult{
		Option:    option,
		Items:     make([]*model.Item, 0, len(data.Records)),
		TotolPage: data.PageCount,
	}

	for _, simple := range data.Records {
		result.Items = append(result.Items, convNhSimpleToItem(simple))
	}
	return result, nil
}

func convNhInfoToItem(info *nhentai.ComicInfo) *model.Item {
	item := &model.Item{
		Title: model.Title{
			Main:       info.Title.English,
			Sub:        info.Title.Japanese,
			Pretty:     info.Title.Pretty,
			MainPretty: utils.GetPrettyTitle(info.Title.English),
			SubPretty:  utils.GetPrettyTitle(info.Title.Japanese),
		},
		Thumb:  nh_rpc.GetThumb(info.MediaId, info.Images.Thumbnail.T),
		Id:     cast.ToString(info.Id),
		Source: consts.NHentai,
		Extra: model.Extra{
			MediaId: info.MediaId,
		},
		UploadDate: info.UploadDate,
		Pages:      info.NumPages,
	}
	extList := make([]string, 0, len(info.Images.Pages))
	for _, page := range info.Images.Pages {
		extList = append(extList, page.T)
	}
	item.Extra.PicExtList = extList
	for _, tag := range info.Tags {
		tagType := strings.ToLower(tag.Type)
		tagValue := tag.Name
		if tagType == consts.TagType_DateAdded || tagType == consts.TagType_Pages {
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

func convNhSimpleToItem(simple nhentai.ComicSimple) *model.Item {
	item := &model.Item{
		Id:     cast.ToString(simple.Id),
		Source: consts.NHentaiSimple,
		Thumb:  simple.Thumb,
		Title: model.Title{
			Main:       simple.Title,
			MainPretty: utils.GetPrettyTitle(simple.Title),
		},
		Lang: simple.Lang,
		Extra: model.Extra{
			MediaId: simple.MediaId,
		},
	}
	return item
}
