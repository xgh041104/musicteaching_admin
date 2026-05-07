package model

type GiveScoreTheme struct {
	GiveScoreThemeId    int    `json:"giveScoreThemeId"`
	GiveScoreCategoryId int    `json:"giveScoreCategoryId"`
	GiveScoreThemeTitle string `json:"giveScoreThemeTitle"`
	CreateTime          string `json:"createTime"`
	SchoolId            int    `json:"schoolId"`
}

type GiveScoreThemeView struct {
	GiveScoreTheme
	GiveScoreCategoryName       string `json:"giveScoreCategoryName"`
	GiveScoreCategoryParentId   int    `json:"giveScoreCategoryParentId"`
	GiveScoreCategoryParentName string `json:"giveScoreCategoryParentName"`
}
