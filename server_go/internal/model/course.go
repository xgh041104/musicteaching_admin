package model

import (
	"time"
)

type Course struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`                // 课程id
	BookID     uint       `gorm:"not null;column:book_id;index" json:"book_id"`      // 新增字段：所属书本 ID
	Title      string     `gorm:"type:varchar(255);not null" json:"title"`           // 课程标题
	VideoPath  string     `gorm:"type:varchar(512)" json:"video_path"`               // 视频路径
	RecordPath string     `gorm:"type:varchar(512)" json:"record_path"`              // 录音路径
	Summary    string     `gorm:"type:text" json:"summary"`                          // 课程总结
	CreatedAt  time.Time  `gorm:"column:create_at" json:"create_at"`                 // 创建时间
	UpdatedAt  time.Time  `gorm:"column:update_at" json:"update_at"`                 // 修改时间
	DeletedAt  *time.Time `gorm:"column:delete_at;index" json:"delete_at,omitempty"` // 删除时间
}

func (m *Course) TableName() string {
	return "course"
}
