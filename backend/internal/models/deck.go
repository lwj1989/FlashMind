package models

import (
	"time"

	"gorm.io/gorm"
)

// Deck 卡包模型
type Deck struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"unique;not null"`
	Archived  bool           `json:"archived" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	Tags  []Tag  `json:"tags,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Cards []Card `json:"cards,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
}

// DeckStats 卡包统计信息
type DeckStats struct {
	TotalCards   int `json:"total_cards"`
	DueCards     int `json:"due_cards"`
	TagCount     int `json:"tag_count"`
	TodayStudied int `json:"today_studied"`
	WeekStudied  int `json:"week_studied"`
}

// DeckWithStats 带统计信息的卡包
type DeckWithStats struct {
	Deck
	Stats DeckStats `json:"stats"`
}
