package model

type ResourcesCategory struct {
	ResourceCategoryId        int    `json:"resourceCategoryId"`
	ResourceCategoryName      string `json:"resourceCategoryName"`
	ResourceCategoryDesc      string `json:"resourceCategoryDesc"`
	ResourceCategoryParentId  int    `json:"resourceCategoryParentId"`
	ResourceCategoryImgFileId int    `json:"resourceCategoryImgFileId"`
	SchoolId                  int    `json:"schoolId"`
	IsDelImg                  int    `json:"isDelImg"`
}

type ResourcesCategoryView struct {
	ResourcesCategory
	FilePath string `json:"filePath"`
	Level    int    `json:"level"`
}

type ResourcesCategoryTreeList struct {
	ResourcesCategoryView
	Children []*ResourcesCategoryTreeList `json:"children"`
}
