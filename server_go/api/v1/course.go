package v1

type CourseRecord struct {
	CourseID uint   `json:"courseId"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Video    string `json:"video"`
	Record   string `json:"record"`
	CreateAt string `json:"createAt"`
}

type CourseListResponse struct {
	Total   int            `json:"total"`
	Records []CourseRecord `json:"records"`
}
type BookRecord struct {
	BookID      uint        `json:"bookId"`
	BookName    string      `json:"bookName"`
	CourseCount int         `json:"courseCount"`
	CreateAt    string      `json:"createAt"`
	UpdateAt    string      `json:"updateAt"`
	DeleteAt    interface{} `json:"deleteAt"` // string 或 null
}

type BookListResponse struct {
	Total   int          `json:"total"`
	Records []BookRecord `json:"records"`
}

type MessageResponse struct {
	Summary string `json:"summary"`
}

type UploadCourseRequest struct {
	Type    int    `form:"type" binding:"required"`
	BookID  uint   `form:"bookId"`
	Title   string `form:"title"  binding:"required"`
	Summary string `form:"summary"` // 非必填
}
