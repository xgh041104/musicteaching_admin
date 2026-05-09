package model

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	UserAccount  string `gorm:"type:varchar(50);not null" json:"userAccount"`
	UserPwd      string `gorm:"type:varchar(255);not null" json:"userPwd"`
	UserType     int    `gorm:"type:int;not null" json:"userType"`
	Remember     bool   `gorm:"type:int;not null" json:"remember"`
	UserTrueName string `gorm:"type:varchar(50);not null" json:"userTrueName"`
	SchoolId     uint   `gorm:"type:int;not null" json:"schoolId"`
}

const (
	UserType = 1
)

func (m *Users) TableName() string {
	return "users"
}
