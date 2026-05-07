package model

type LoginUser struct {
	UserId       int    `json:"userId" `
	UserTrueName string `json:"userTrueName" `
	UserAccount  string `json:"userAccount"  `
	UserPwd      string `json:"userPwd" `
	SchoolId     int    `json:"schoolId"  `
	UserType     int    `json:"userType"  `
	ClassName    string `json:"className"`
	AccessToken  string `json:"accessToken"`
}

type LoginReq struct {
	UserAccount string `json:"userAccount" binding:"required"`
	UserPwd     string `json:"userPwd" binding:"required"`
	UserType    int    `json:"userType" `
	SerialNum   string `json:"serialNum"`
	HostMAC     string `json:"hostMAC"`
	AccessToken string `json:"accessToken"`
	TeacherId   int    `json:"teacherId"`
}

type CommonUser struct {
	CommonUserId       int    `json:"userId" `
	CommonUserTrueName string `json:"userTrueName"  `
	CommonUserAccount  string `json:"userAccount" `
	CommonUserPwd      string `json:"userPwd"  `
	SchoolId           int    `json:"schoolId"  `
	CommonUserType     int    `json:"userType"  `
}
