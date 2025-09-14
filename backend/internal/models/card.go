package models

import (
	"time"

	"gorm.io/gorm"
)

// Card 卡片模型
type Card struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	DeckID    uint           `json:"deck_id" gorm:"not null;index"`
	TagID     *uint          `json:"tag_id,omitempty" gorm:"index"` // 可为空，表示未分组
	Question  string         `json:"question" gorm:"not null;type:text"`
	Answer    string         `json:"answer" gorm:"not null;type:text"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	Deck   Deck    `json:"deck,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Tag    *Tag    `json:"tag,omitempty" gorm:"constraint:OnDelete:SET NULL;"`
	Review *Review `json:"review,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
}

// CardSearchRequest 卡片搜索请求
type CardSearchRequest struct {
	DeckID   *uint  `form:"deck_id"`
	TagID    *uint  `form:"tag_id"`
	Keyword  string `form:"keyword"`
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
}

// CardResponse 卡片响应（包含关联数据）
type CardResponse struct {
	Card
	DeckName string `json:"deck_name,omitempty"`
	TagName  string `json:"tag_name,omitempty"`
}

// CardListResponse 卡片列表响应
type CardListResponse struct {
	Cards      []CardResponse `json:"cards"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}
