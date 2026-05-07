package model

type Admin struct {
	AdminId       int    `json:"userId" db:"AdminId"`
	AdminTrueName string `json:"userTrueName" db:"AdminTrueName"`
	AdminAccount  string `json:"userAccount" db:"AdminAccount"`
	AdminPwd      string `json:"userPwd" db:"AdminPwd"`
	SchoolId      int    `json:"schoolId" db:"SchoolId"`
	AdminType     int    `json:"userType" db:"AdminType"`
	SchoolName    string `json:"schoolName" db:"SchoolName"`
}
