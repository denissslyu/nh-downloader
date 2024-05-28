package nh_serv

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cast"

	"nh-downloader/consts"
	"nh-downloader/internal/config"
	"nh-downloader/internal/model"
	"nh-downloader/internal/rpc/nh_rpc"
	"nh-downloader/utils"
	"nh-downloader/utils/logs"
)

type infoJson struct {
	GalleryInfo struct {
		Title              string `json:"title"`
		TitleTitleOriginal string `json:"title_title_original"`
		Link               string `json:"link"`
		Category           string `json:"category,omitempty"`
		Tags               struct {
			Group    []string `json:"group,omitempty"`
			Language []string `json:"language,omitempty"`
			Parody   []string `json:"parody,omitempty"`
			Artist   []string `json:"artist,omitempty"`
			Tag      []string `json:"tag,omitempty"`
			Category []string `json:"category,omitempty"`
		} `json:"tags"`
		Language   string `json:"language"`
		Translated bool   `json:"translated"`
		UploadDate []int  `json:"upload_date"`
		Source     struct {
			Site string `json:"site"`
			Gid  int    `json:"gid"`
		} `json:"source"`
	} `json:"gallery_info"`
}

func downloadPage(ext, folder string, mediaId, page int) error {
	var err error
	logs.Info("[nh_serv.downloadPage] starting to downloadPage: ", mediaId, page)
	extension := nh_rpc.Extension(ext)
	saveFilePath := path.Join(folder, fmt.Sprintf("%d.%s", page, extension))
	if _, err = os.Stat(saveFilePath); err == nil {
		logs.Info("[nh_serv.downloadPage] ignore saved file:", saveFilePath)
		return nil
	}

	file, err := os.Create(saveFilePath)
	if err != nil {
		logs.Error("[nh_serv.downloadPage] create file failed: ", err)
		return err
	}
	defer file.Close()

	for i := 0; i < config.Retried(); i++ {
		resp, er := nh_rpc.Get(nh_rpc.PageUrl(mediaId, page, ext))
		if er != nil {
			err = er
			logs.Warn("[nh_serv.downloadPage] request failed: ", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			err = errors.New("statusCode not 200")
			logs.Warn("[nh_serv.downloadPage] StatusCode invalid: ", resp.StatusCode)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		_, er = io.Copy(file, resp.Body)
		if er != nil {
			err = er
			logs.Warn("[nh_serv.downloadPage] file copy failed:", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if er == nil {
			err = er
			break
		}
	}
	if err != nil {
		logs.Error("[nh_serv.downloadPage] download page,retried failed")
		return err
	}

	return nil
}

func DownloadByItemId(id string) error {
	item, err := GetItem(id)
	if err != nil {
		return err
	}
	return DownloadItem(item)
}

func DownloadItem(item *model.Item) error {
	var err error
	// prepare folder
	folder := path.Join(config.CachePath(), strings.ReplaceAll(item.Title.Main, "/", ""))
	if _, err = os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			logs.Error("[nh_serv.DownloadItem] mkdir failed:", err)
			return err
		}
	} else {
		logs.Info("[nh_serv.DownloadItem] folder exist,remove files")
		err = deleteAllFiles(folder)
		if err != nil {
			logs.Error("[nh_serv.DownloadItem] remove files failed,", err)
			return err
		}
	}

	// remove cache folder
	defer func() {
		err = os.RemoveAll(folder)
		if err != nil {
			logs.Error("[nh_serv.DownloadItem] remove cache failed:", err)
		}
	}()

	// todo config
	maxGoroutines := 5 // 最大 goroutine 数量
	goroutineCh := make(chan int, maxGoroutines)
	errCh := make(chan error, 1)
	wg := sync.WaitGroup{}

	// download every page of item
	for idx, ext := range item.Extra.PicExtList {
		page := idx + 1
		wg.Add(1)
		goroutineCh <- page

		go func(ext string, folder string, mediaID int, page int) {
			defer wg.Done()

			err = downloadPage(ext, folder, mediaID, page)
			if err != nil {
				errCh <- err
			} else {
				<-goroutineCh
			}
		}(ext, folder, item.Extra.MediaId, page)
	}

	// wait for all pages downloaded
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// download eze style info.json
	json := makeEzeInfoJson(item)
	err = utils.SaveJsonToFile(path.Join(folder, "info.json"), json)
	if err != nil {
		logs.Error("[nh_serv.DownloadItem] download info.json failed:", err)
		return err
	}

	for err = range errCh {
		if err != nil {
			logs.Error("[nh_serv.DownloadItem] download failed:", err)
			return err
		}
	}

	// zip file
	err = utils.ZipFiles(folder, path.Join(config.DownloadPath(), strings.ReplaceAll(item.Title.Main+".zip", "/", "")))
	if err != nil {
		logs.Error("[nh_serv.DownloadItem] zip failed:", err)
		return err
	}

	return nil
}

func deleteAllFiles(directory string) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			err := os.Remove(path)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func makeEzeInfoJson(item *model.Item) *infoJson {
	info := &infoJson{}
	info.GalleryInfo.Title = item.Title.Main
	info.GalleryInfo.TitleTitleOriginal = item.Title.Sub
	info.GalleryInfo.Link = consts.NhentaiGalleryUrlPrefix + item.Id
	info.GalleryInfo.Category = item.GetCategory()
	info.GalleryInfo.Language = item.GetLanguage()
	info.GalleryInfo.Translated = item.IsTranslated()
	info.GalleryInfo.Source.Gid = cast.ToInt(item.Id)
	info.GalleryInfo.Source.Site = "https://nhentai.net"

	updateDate := int(time.Now().Unix())
	if item.UploadDate > 0 {
		updateDate = item.UploadDate
	}
	t := time.Unix(int64(updateDate), 0)
	info.GalleryInfo.UploadDate = []int{
		t.Year(),
		int(t.Month()),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
	}

	for _, tag := range item.Tags {
		switch tag.Type {
		case consts.TagType_Artist:
			info.GalleryInfo.Tags.Artist = append(info.GalleryInfo.Tags.Artist, tag.Value)
		case consts.TagType_Language:
			info.GalleryInfo.Tags.Language = append(info.GalleryInfo.Tags.Language, tag.Value)
		case consts.TagType_Category:
			info.GalleryInfo.Tags.Category = append(info.GalleryInfo.Tags.Category, tag.Value)
		case consts.TagType_Group:
			info.GalleryInfo.Tags.Group = append(info.GalleryInfo.Tags.Group, tag.Value)
		case consts.TagType_Parody:
			info.GalleryInfo.Tags.Parody = append(info.GalleryInfo.Tags.Parody, tag.Value)
		default:
			info.GalleryInfo.Tags.Tag = append(info.GalleryInfo.Tags.Tag, tag.Value)
		}
	}
	return info
}
