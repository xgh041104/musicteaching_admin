package model

type Course struct {
	CourseId             int    `json:"courseId" `
	CourseCategoryId     int    `json:"courseCategoryId" `
	CourseTitle          string `json:"courseTitle" `
	CourseDesc           string `json:"courseDesc" `
	SchoolId             int    `json:"schoolId" `
	LecturerCommonUserId int    `json:"lecturerCommonUserId" `
	CourseUpdateTime     string `json:"courseUpdateTime" `
	CourseImgFileId      int    `json:"courseImgFileId" `
	ViewNum              int    `json:"viewNum" `
	LikeNum              int    `json:"likeNum" `
	IsDelImg             int    `json:"isDelImg" `
	CourseType           int    `json:"courseType" `
	Status               int    `json:"status"`
}

type CourseView struct {
	Course
	IsCollected              bool   `json:"isCollected"`
	FilePath                 string `json:"filePath" `
	CourseCategoryName       string `json:"courseCategoryName" `
	SchoolName               string `json:"schoolName" `
	CommonUserTrueName       string `json:"commonUserTrueName" `
	CourseCategoryParentName string `json:"courseCategoryParentName" `
	CourseCategoryParentId   int    `json:"courseCategoryParentId" `
}

type CourseDirectory struct {
	CourseView
	ChapterArr []*ChapterDirectory `json:"chapterArr" `
}

type CourseResult struct {
	Total   int           `json:"total"`
	Records []*CourseView `json:"records"`
}
