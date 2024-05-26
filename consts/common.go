package consts

const (
	NHentai       = "nhentai"
	NHentaiSimple = "nhentai_simple"
	Lanraragi     = "lanraragi"
)

const (
	TagType_Artist     = "artist"
	Tagtype_Artists    = "artists"
	TagType_Group      = "group"
	TagType_Groups     = "groups"
	TagType_Parody     = "parody"
	TagType_Parodies   = "parodies"
	TagType_Category   = "category"
	TagType_Categories = "categories"
	TagType_Language   = "language"
	TagType_Languages  = "languages"
	TagType_Tag        = "tag"
	TagType_Tags       = "tags"
	TagType_Character  = "character"
	TagType_DateAdded  = "date_added"
	TagType_Male       = "male"
	TagType_Female     = "female"
	TagType_Other      = "other"
	TagType_Mixed      = "mixed"
	TagType_Source     = "source"
	TagType_Pages      = "pages"

	LanguageTag_Translated = "translated"
	LanguageTag_Japanese   = "japanese"
	LanguageTag_English    = "english"
	LanguageTag_Chinese    = "chinese"
)

var TagTypeReplaceMap = map[string]string{
	Tagtype_Artists:    TagType_Artist,
	TagType_Groups:     TagType_Group,
	TagType_Parodies:   TagType_Parody,
	TagType_Categories: TagType_Category,
	TagType_Languages:  TagType_Language,
	TagType_Tags:       TagType_Tag,
}

var TagTypes = []string{
	TagType_Artist,
	TagType_Group,
	TagType_Parody,
	TagType_Category,
	TagType_Language,
	TagType_Character,
	TagType_DateAdded,
	TagType_Male,
	TagType_Female,
	TagType_Tag,
	TagType_Other,
	TagType_Mixed,
	TagType_Source,
}

const (
	Lang_CH = "CH"
	Lang_JP = "JP"
	Lang_EN = "EN"
)

const (
	NhPageSize  = 25
	LrrPageSize = 100
)
