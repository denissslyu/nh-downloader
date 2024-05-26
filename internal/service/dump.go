package service

import (
	"bytes"
	"image"
	"image/jpeg"
	"net/http"
	"nh-downloader/internal/config"
	"nh-downloader/internal/rpc/lrr_rpc"
	"strings"
	"time"

	"nh-downloader/consts"
	"nh-downloader/internal/model"
	"nh-downloader/internal/service/lrr_serv"
	"nh-downloader/internal/service/nh_serv"
	"nh-downloader/utils/logs"

	"github.com/corona10/goimagehash"
	"github.com/gojp/kana"
)

type dumpContext struct {
	LrrItemMap              map[string]*model.Item
	CachedItemMap           map[string]*model.Item
	DownloadedItemMap       map[string]*model.Item
	DownloadFailItemMap     map[string]*model.Item
	IdemItemInfoMap         map[string]*idemItemInfo         //key: nhId
	SimilarThumbItemInfoMap map[string]*similarThumbItemInfo // key: nhId
}

type idemItemInfo struct {
	Item     *model.Item
	IdemType int
}

type similarThumbItemInfo struct {
	Item  *model.Item
	LrrId string
	Score float64
}

// Dump to dump items by simple search
func Dump(filters []string) error {
	ctx := &dumpContext{
		LrrItemMap:              make(map[string]*model.Item),
		CachedItemMap:           make(map[string]*model.Item),
		DownloadedItemMap:       make(map[string]*model.Item),
		DownloadFailItemMap:     make(map[string]*model.Item),
		IdemItemInfoMap:         make(map[string]*idemItemInfo),
		SimilarThumbItemInfoMap: make(map[string]*similarThumbItemInfo),
	}
	// get filtered items from lanraragi
	lrrResult, err := lrr_serv.SimpleSearchItems(&model.SimpleSearchOption{
		Filters: filters,
		Page:    -1,
	})
	if err != nil {
		logs.Error("[service.Dump] search lrr items failed:", err)
		return err
	}
	if lrrResult != nil {
		for _, lrrItem := range lrrResult.Items {
			ctx.LrrItemMap[lrrItem.Id] = lrrItem
		}
		logs.Info("[service.Dump] filtered lrr items,len:", len(ctx.LrrItemMap))
	}
	page := 1
	var total int
	for {
		total = dumpWithPage(filters, page, ctx)
		page++
		if total < page {
			break
		}
	}

	// find similar thumbnail
	err = findSimilarThumbnail(ctx)
	if err != nil {
		logs.Error("[service.Dump] find Similar Thumbnail failed:", err)
		return err
	}

	for itemId, item := range ctx.CachedItemMap {
		err = nh_serv.DownloadItem(item)
		if err != nil {
			ctx.DownloadFailItemMap[itemId] = item
		} else {
			ctx.DownloadedItemMap[itemId] = item
		}
	}

	// report
	logs.Info("[service.Dump] -------------------------⬇ idempotent items ⬇-------------------------")
	for itemId, itemInfo := range ctx.IdemItemInfoMap {
		logs.Info("[service.Dump]【url】 https://nhentai.net/g/"+itemId, "【idemType】", itemInfo.IdemType, "【title】", itemInfo.Item.Title.Main)
	}
	logs.Info("[service.Dump] ----------------------⬇ similar thumb items ⬇----------------------")
	for itemId, itemInfo := range ctx.SimilarThumbItemInfoMap {
		logs.Info("[service.Dump]【url】 https://nhentai.net/g/"+itemId, "【lrrId】", itemInfo.LrrId, "【title】", itemInfo.Item.Title.Main, "【score】", itemInfo.Score)
	}
	logs.Info("[service.Dump] -------------------------⬇ dumped items ⬇-------------------------")
	for itemId, item := range ctx.DownloadedItemMap {
		logs.Info("[service.Dump] dump success 【url】 https://nhentai.net/g/"+itemId, "【title】", item.Title.Main)
	}
	logs.Info("[service.Dump] -------------------------⬇ dumped failed items ⬇-------------------------")
	for itemId, item := range ctx.DownloadFailItemMap {
		logs.Error("[service.Dump] dump failed 【url】 https://nhentai.net/g/"+itemId, "【title】", item.Title.Main)
	}
	return nil
}

func dumpWithPage(filters []string, page int, ctx *dumpContext) int {
	option := &model.SimpleSearchOption{
		Filters: filters,
		Page:    page,
	}
	var result *model.SimpleSearchResult
	var err error
	// retry
	for i := 0; i < config.Retried(); i++ {
		result, err = nh_serv.SimpleSearchItems(option)
		if err == nil {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		logs.Error("[service.dumpWithPage] retried search nhItems failed")
		return -1
	}
	if result == nil {
		logs.Warn("[service.dumpWithPage] search nhItems empty")
		return -1
	}
	logs.Info("[service.Dump] start dump page", page, ", total", result.TotolPage)
	for _, item := range result.Items {
		if item.Lang != consts.Lang_CH {
			continue
		}

		// check idempotent: by mainpretty --- cached items
		if idemItem := idempotentByTitle(item.Title.MainPretty+"|"+item.Title.SubPretty, ctx.CachedItemMap); idemItem != nil {
			ctx.IdemItemInfoMap[item.Id] = &idemItemInfo{Item: item, IdemType: consts.DumpIdemTypeCached}
			continue
		}

		// check idempotent: by mainpretty --- lrr items
		if idemItem := idempotentByTitle(item.Title.MainPretty, ctx.LrrItemMap); idemItem != nil {
			ctx.IdemItemInfoMap[item.Id] = &idemItemInfo{Item: item, IdemType: consts.DumpIdemTypeLrr}
			continue
		}

		// check idempotent: by subpretty --- lrr items
		for i := 0; i < config.Retried(); i++ {
			item, err = nh_serv.GetItem(item.Id)
			if err == nil {
				break
			}

			time.Sleep(500 * time.Millisecond)
		}
		if err != nil {
			logs.Error("[service.dumpWithPage] retried getItem failed")
			return -1
		}
		if idemItem := idempotentByTitle(item.Title.SubPretty, ctx.LrrItemMap); idemItem != nil {
			ctx.IdemItemInfoMap[item.Id] = &idemItemInfo{Item: item, IdemType: consts.DumpIdemTypeLrr}
			continue
		}

		ctx.CachedItemMap[item.Id] = item
	}
	return result.TotolPage
}

func idempotentByTitle(title string, itemMap map[string]*model.Item) *model.Item {
	title = kana.KanaToRomaji(strings.ToLower(strings.ReplaceAll(title, " ", "")))
	for _, subTitle := range strings.Split(title, "|") {
		if subTitle == "" {
			continue
		}
		for _, item := range itemMap {
			combinedPretty := item.Title.MainPretty + "|" + item.Title.SubPretty + "|" + item.Title.Pretty
			combinedPretty = kana.KanaToRomaji(strings.ToLower(strings.ReplaceAll(combinedPretty, " ", "")))
			for _, pretty := range strings.Split(combinedPretty, "|") {
				if pretty == "" {
					continue
				}
				if subTitle == pretty {
					return item
				}
			}
		}
	}
	return nil
}

func loadImageFromBytes(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		logs.Error("[service.loadImageFromBytes] decode image failed:", err)
		return nil, err
	}
	return img, nil
}

// 从URL加载图片
func loadImageFromUrl(url string) (image.Image, error) {
	var resp *http.Response
	var err error

	for i := 0; i < config.Retried(); i++ {
		resp, err = http.Get(url)
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		logs.Error("[service.loadImageFromUrl] http get retried failed:", err)
		return nil, err
	}
	defer resp.Body.Close()
	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		logs.Error("[service.loadImageFromUrl] decode image failed:", err)
		return nil, err
	}
	return img, nil
}

// similarItemMap: similar thumbnail item in lrr and nh, key: nhId, val: lrrId
func findSimilarThumbnail(ctx *dumpContext) error {
	// prepare lrr thumb hash map
	lrrThumbHashMap := make(map[string]*goimagehash.ImageHash)
	for _, item := range ctx.LrrItemMap {
		bytes, err := lrr_rpc.GetThumbBytes(item.Id)
		if err != nil {
			return err
		}
		img, err := loadImageFromBytes(bytes)
		if err != nil {
			return err
		}
		hash, err := goimagehash.PerceptionHash(img)
		if err != nil {
			logs.Error("[service.checkSimilarThumbnail] perceptionHash failed:", err)
			return err
		}
		lrrThumbHashMap[item.Id] = hash
	}
	// todo config
	waterLine := 0.8
	nhIdsTobeDelFromCached := make([]string, 0)
	for itemId, item := range ctx.CachedItemMap {
		logs.Info("[service.checkSimilarThumbnail] check Similarity of ", item.Title.Main)
		img, err := loadImageFromUrl(item.Thumb)
		if err != nil {
			continue
		}
		hash, err := goimagehash.PerceptionHash(img)
		if err != nil {
			logs.Error("[service.checkSimilarThumbnail] perceptionHash failed:", err)
			continue
		}
		score := 0.0
		lrrId := ""
		for id, lrrHash := range lrrThumbHashMap {
			similarity, err := calculateSimilarity(hash, lrrHash)
			if err != nil {
				continue
			}

			if similarity > score {
				score = similarity
				lrrId = id
			}
		}
		if score >= waterLine {
			nhIdsTobeDelFromCached = append(nhIdsTobeDelFromCached, itemId)
			ctx.SimilarThumbItemInfoMap[itemId] = &similarThumbItemInfo{Item: item, LrrId: lrrId, Score: score}
		}
	}
	for _, id := range nhIdsTobeDelFromCached {
		delete(ctx.CachedItemMap, id)
	}
	return nil
}

// 计算图片相似度
func calculateSimilarity(hashA, hashB *goimagehash.ImageHash) (float64, error) {
	distance, err := hashA.Distance(hashB)
	if err != nil {
		logs.Error("[service.calculateSimilarity] calculate distance failed,", err)
		return 0, err
	}
	maxDistance := 64 // PerceptionHash generates a 64-bit hash
	similarity := 1 - float64(distance)/float64(maxDistance)
	return similarity, nil
}
