package model

type CourseCategory struct {
	CourseCategoryId        int    `json:"courseCategoryId"`
	CourseCategoryName      string `json:"courseCategoryName"`
	CourseCategoryParentId  int    `json:"courseCategoryParentId"`
	SchoolId                int    `json:"schoolId"`
	CourseCategoryImgFileId int    `json:"courseCategoryImgFileId"`
	IsDelImg                int    `json:"isDelImg" `
}

type CourseCategoryView struct {
	CourseCategory
	FilePath string `json:"filePath"`
	Level    int    `json:"level"`
}

type CourseCategoryTreeList struct {
	CourseCategoryView
	Children []*CourseCategoryTreeList `json:"children"`
}
