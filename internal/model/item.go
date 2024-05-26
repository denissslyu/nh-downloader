package model

import "nh-downloader/consts"

type Item struct {
	Id         string
	Thumb      string
	Source     string
	Lang       string
	Title      Title
	Tags       []Tag
	Extra      Extra
	Pages      int
	UploadDate int
}

func (i *Item) SetLangByTag() {
	for _, tag := range i.Tags {
		if tag.Type == consts.TagType_Language {
			switch tag.Value {
			case consts.LanguageTag_Chinese:
				i.Lang = consts.Lang_CH
			case consts.LanguageTag_Japanese:
				i.Lang = consts.Lang_JP
			case consts.LanguageTag_English:
				i.Lang = consts.Lang_EN
			}
		}
	}
}

func (i *Item) GetLanguage() string {
	for _, tag := range i.Tags {
		if tag.Type == consts.TagType_Language && tag.Value != consts.LanguageTag_Translated {
			return tag.Value
		}
	}
	return ""
}

func (i *Item) IsTranslated() bool {
	for _, tag := range i.Tags {
		if tag.Type == consts.TagType_Language && tag.Value == consts.LanguageTag_Translated {
			return true
		}
	}
	return false
}

func (i *Item) GetCategory() string {
	for _, tag := range i.Tags {
		if tag.Type == consts.TagType_Category {
			return tag.Value
		}
	}
	return ""
}

type Extra struct {
	// MediaId for nhentai
	MediaId int
	// PicExtList for nhentai
	PicExtList []string
}

type Title struct {
	// EnglishTitle of nhentai; Title of lanraragi
	Main string
	// JapaneseTitle of nhentai
	Sub string
	// Pretty of nhentai
	Pretty string
	// Pretty by Main
	MainPretty string
	// Pretty by Sub
	SubPretty string
}

type Tag struct {
	Type  string
	Value string
}

type SimpleSearchOption struct {
	// filter words for searching
	Filters []string
	// both lanraragi,nhentai are start from 1
	// pagesize cannot change,25 for nhentai, 100 for lanraragi
	Page int
	// sort
	// nhentai-go doesn't support, maybe commit a PR in future
	// lanraragi: title, data, artist, group, pages, language, category, female, male, other, tag
	Sort string
	// only for lanraragi, support asc or desc
	Order string
}

type SimpleSearchResult struct {
	Items     []*Item
	Option    *SimpleSearchOption
	TotolPage int
}
