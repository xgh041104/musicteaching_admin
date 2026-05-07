package model

type StudentUser struct {
	StudentId      int    `json:"studentId"`
	StudentAccount string `json:"studentAccount"`
	StudentPwd     string `json:"studentPwd"`
	StudentName    string `json:"studentName"`
	ClassName      string `json:"className"`
	SchoolId       int    `json:"schoolId"`
	UserType       int    `json:"userType"`
}

type Student struct {
	StudentUser
	TeacherIdList []int `json:"teacherIdList"`
}
