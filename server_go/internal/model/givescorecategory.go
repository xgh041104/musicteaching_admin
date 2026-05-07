package model

type GiveScoreCategory struct {
	GiveScoreCategoryId        int    `json:"giveScoreCategoryId"`
	GiveScoreCategoryName      string `json:"giveScoreCategoryName"`
	GiveScoreCategoryDesc      string `json:"giveScoreCategoryDesc"`
	GiveScoreCategoryParentId  int    `json:"giveScoreCategoryParentId"`
	SchoolId                   int    `json:"schoolId"`
	GiveScoreCategoryImgFileId int    `json:"giveScoreCategoryImgFileId"`
	IsDelImg                   int    `json:"isDelImg" `
}

type GiveScoreCategoryView struct {
	GiveScoreCategory
	FilePath string `json:"filePath"`
	Level    int    `json:"level"`
}

type GiveScoreCategoryTreeList struct {
	GiveScoreCategoryView
	Children []*GiveScoreCategoryTreeList `json:"children"`
}
