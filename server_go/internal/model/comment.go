package model

type Comment struct {
	CommentId      int    `json:"commentId"`
	CourseId       int    `json:"courseId"`
	ChapterId      int    `json:"chapterId"`
	SectionId      int    `json:"sectionId"`
	CommentContent string `json:"commentContent"`
	CommonUserId   int    `json:"commonUserId"`
	CommentTime    string `json:"commentTime"`
}

type CommenntView struct {
	Comment
	CommonUserTrueName string `json:"commonUserTrueName"`
	CourseTitle        string `json:"courseTitle"`
	ChapterTitle       string `json:"chapterTitle"`
	SectionTitle       string `json:"sectionTitle"`
}
