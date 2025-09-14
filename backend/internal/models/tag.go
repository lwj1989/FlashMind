package models

import (
	"time"

	"gorm.io/gorm"
)

// Tag 标签模型
type Tag struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	DeckID    *uint          `json:"deck_id,omitempty" gorm:"index"`
	Name      string         `json:"name" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	Deck  *Deck  `json:"deck,omitempty" gorm:"constraint:OnDelete:SET NULL;"`
	Cards []Card `json:"cards,omitempty"`
}

// TagStats 标签统计信息
type TagStats struct {
	TotalCards   int `json:"total_cards"`
	DueCards     int `json:"due_cards"`
	TodayStudied int `json:"today_studied"`
}

// TagWithStats 带统计信息的标签
type TagWithStats struct {
	Tag
	Stats TagStats `json:"stats"`
}

// 确保同一卡包内标签名唯一的索引
func (Tag) TableName() string {
	return "tags"
}
