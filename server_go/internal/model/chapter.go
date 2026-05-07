package model

type Chapter struct {
	ChapterId    int    `json:"chapterId"  `
	CourseId     int    `json:"courseId" `
	ChapterTitle string `json:"chapterTitle" `
	ChapterOrder int    `json:"chapterOrder" `
}
type ChapterView struct {
	Chapter
	CourseTitle string `json:"courseTitle"  `
}

type ChapterDirectory struct {
	Chapter
	SectionArr []*Section `json:"sectionArr"  `
}
