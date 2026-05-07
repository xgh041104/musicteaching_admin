package model

type Resources struct {
	ResourceId           int    `json:"resourceId"`
	ResourceCategoryId   int    `json:"resourceCategoryId" `
	ResourceName         string `json:"resourceName" `
	ResourceType         string `json:"resourceType"`
	ResourceFileId       int    `json:"resourceFileId"`
	SchoolId             int    `json:"schoolId"`
	ResourceDesc         string `json:"resourceDesc"`
	ResourceIsPublic     int    `json:"resourceIsPublic"`
	ResourceImgFileId    int    `json:"resourceImgFileId"`
	LecturerCommonUserId int    `json:"lecturerCommonUserId"`
}

type ResourcesView struct {
	Resources
	ResourceFilePath     string `json:"resourceFilePath"`
	ResourceFileName     string `json:"resourceFileName"`
	ImgfilePath          string `json:"imgfilePath"`
	ResourceCategoryName string `json:"resourceCategoryName"`
	SchoolName           string `json:"schoolName"`
	CommonUserTrueName   string `json:"commonUserTrueName"`
}
