package model

type GiveScore struct {
	GeScoreId      int     `json:"giveScoreId"`
	GeScoreThemeId int     `json:"giveScoreThemeId"`
	StudentName    string  `json:"studentName"`
	StudentGrade   string  `json:"studentGrade"`
	StudentSubject string  `json:"studentSubject"`
	TeacherScore   float32 `json:"teacherScore"`
	OneselfScore   float32 `json:"oneselfScore"`
	MutualScore    float32 `json:"mutualScore"`
}
