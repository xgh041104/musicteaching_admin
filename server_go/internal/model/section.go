package model

type Section struct {
	SectionId    int    `json:"sectionId"`
	ChapterId    int    `json:"chapterId"`
	SectionTitle string `json:"sectionTitle"`
	SectionDesc  string `json:"sectionDesc"`
	//0 图文课 1视频课 2图片课 3 ppt课
	SectionType int `json:"sectionType"`

	SectionContent string `json:"sectionContent"`
	SectionOrder   int    `json:"sectionOrder"`
}

type SectionrelationFile struct {
	SectionrelationId   int `json:"sectionrelationId"`
	SectionId           int `json:"sectionId"`
	FileInfoId          int `json:"fileInfoId"`
	SectionFileOrder    int `json:"sectionFileOrder"`
	SectionFileInfoType int `json:"sectionFileInfoType"`
}

type SectionView struct {
	Section
	ChapterTitle string                     `json:"chapterTitle"`
	FileContent  []*SectionrelationFileView `json:"fileContent"`
	FileAnnex    []*SectionrelationFileView `json:"fileAnnex"`
}

type SectionModel struct {
	Section
	ChapterTitle string
}
type SectionrelationFileView struct {
	SectionrelationId   int `json:"sectionrelationId"`
	SectionFileInfoType int `json:"sectionFileInfoType"`
	SectionFileOrder    int `json:"sectionFileOrder"`
	FileInfo
}

type SectionEdit struct {
	Section
	RemoveFile       []int          `json:"removeFile" `
	NewResourceFiles []*ResourceAdd `json:"newResourceFiles" `
}

type ResourceAdd struct {
	ResourceId   int `json:"resourceId"`
	Resourcetype int `json:"resourcetype"`
	Position     int `json:"position"`
}

// type SectionView struct {
// 	Section
// }
