package models

import "time"

// BackupData 完整备份数据结构
type BackupData struct {
	Version    string         `json:"version"`
	ExportDate time.Time      `json:"export_date"`
	Decks      []DeckBackup   `json:"decks"`
	Tags       []TagBackup    `json:"tags"`
	Cards      []CardBackup   `json:"cards"`
	Reviews    []ReviewBackup `json:"reviews"`
}

// 完整表备份结构
type DeckBackup struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Archived  bool      `json:"archived"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TagBackup struct {
	ID        uint      `json:"id"`
	DeckID    *uint     `json:"deck_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CardBackup struct {
	ID        uint      `json:"id"`
	DeckID    uint      `json:"deck_id"`
	TagID     *uint     `json:"tag_id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReviewBackup struct {
	ID          uint      `json:"id"`
	CardID      uint      `json:"card_id"`
	EFactor     float64   `json:"efactor"`
	Interval    int       `json:"interval"`
	Repetitions int       `json:"repetitions"`
	NextReview  time.Time `json:"next_review"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 为了兼容性，保留原有导出结构
type DeckExport struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Cards       []CardExport `json:"cards"`
	Tags        []TagExport  `json:"tags"`
	CreatedAt   time.Time    `json:"created_at"`
}

type CardExport struct {
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	TagName   string    `json:"tag_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type TagExport struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// SystemStats 系统统计信息
type SystemStats struct {
	TotalDecks int    `json:"total_decks"`
	TotalCards int    `json:"total_cards"`
	TotalTags  int    `json:"total_tags"`
	Version    string `json:"version"`
}
