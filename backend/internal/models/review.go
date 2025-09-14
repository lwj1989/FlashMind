package models

import (
	"time"
)

// Review 复习调度模型 (SM-2算法相关)
type Review struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CardID      uint      `json:"card_id" gorm:"unique;not null;index"` // 1:1关系
	EFactor     float64   `json:"efactor" gorm:"default:2.5"`           // 记忆强度因子，默认2.5
	Interval    int       `json:"interval" gorm:"default:0"`            // 间隔天数
	Repetitions int       `json:"repetitions" gorm:"default:0"`         // 连续复习次数
	NextReview  time.Time `json:"next_review" gorm:"index"`             // 下次复习时间
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	Card Card `json:"card,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
}

// ReviewResult 复习结果枚举
type ReviewResult int

const (
	Again ReviewResult = iota // 忘记，重新学习
	Hard                      // 模糊，稍微难记
	Good                      // 记得，正常难度
)

// StudyQueue 学习队列项
type StudyQueue struct {
	CardID   uint   `json:"card_id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	DeckName string `json:"deck_name"`
	TagName  string `json:"tag_name,omitempty"`
}

// StudySession 学习会话
type StudySession struct {
	Queue     []StudyQueue `json:"queue"`
	Current   int          `json:"current"`
	Total     int          `json:"total"`
	Completed int          `json:"completed"`
	StartTime time.Time    `json:"start_time"`
}

// ReviewRequest 复习请求
type ReviewRequest struct {
	Result ReviewResult `json:"result"`
}

// ReviewResponse 复习响应
type ReviewResponse struct {
	Success    bool      `json:"success"`
	NextReview time.Time `json:"next_review"`
	Interval   int       `json:"interval"`
	Message    string    `json:"message"`
}
