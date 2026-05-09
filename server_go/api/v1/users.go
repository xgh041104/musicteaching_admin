package v1

type Login struct {
	UserAccount string `json:"userAccount" binding:"required"`
	UserPwd     string `json:"userPwd" binding:"required"`
	UserType    int    `json:"userType" binding:"required"`
	Remember    bool   `json:"remember"`
}
type LoginResponse struct {
	Data  *User  `json:"data"`
	Token string `json:"token"`
}
type User struct {
	ID           uint   `json:"userId"`
	SchoolId     uint   `json:"schoolId"`
	UserAccount  string `json:"userAccount"`
	UserPwd      string `json:"userPwd"`
	UserTrueName string `json:"userTrueName"`
	UserType     int    `json:"userType"`
}
