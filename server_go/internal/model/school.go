package model

type School struct {
	SchoolId        int    `json:"schoolId"`
	SchoolName      string `json:"schoolName"`
	Address         string `json:"address"`
	SchoolImgFileId int    `json:"schoolImgFileId"`
	IsDelImg        int    `json:"isDelImg" `
	PaidCourseIds   []int  `json:"paidCourseIds"`
}

type SchoolView struct {
	SchoolId          int      `json:"schoolId"`
	SchoolName        string   `json:"schoolName"`
	Address           string   `json:"address"`
	SchoolImgFileId   int      `json:"schoolImgFileId"`
	UserNum           int      `json:"userNum"`
	SchoolImgFileName string   `json:"schoolImgFileName"`
	SchoolImgPath     string   `json:"schoolImgPath"`
	PaidCourseIds     []int    `json:"paidCourseIds"`
	PaidCourseTitles  []string `json:"paidCourseTitles"`
}
