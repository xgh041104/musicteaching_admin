package model

import "time"

type Book struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`                // 书本Id，主键
	BookName    string     `gorm:"type:varchar(255);not null" json:"book_name"`       // 书本名称
	CourseCount int        `gorm:"not null" json:"course_count"`                      // 视频数
	CreatedAt   time.Time  `gorm:"column:create_at" json:"create_at"`                 // 创建时间
	UpdatedAt   time.Time  `gorm:"column:update_at" json:"update_at"`                 // 修改时间
	DeletedAt   *time.Time `gorm:"column:delete_at;index" json:"delete_at,omitempty"` // 删除时间（软删除）
}

func (Book) TableName() string {
	return "book"
}
