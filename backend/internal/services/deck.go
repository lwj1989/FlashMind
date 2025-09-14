package services

import (
	"flashcard/internal/models"
	"flashcard/pkg/database"

	"gorm.io/gorm"
)

// DeckService 卡包服务
type DeckService struct {
	db *gorm.DB
}

// NewDeckService 创建卡包服务实例
func NewDeckService() *DeckService {
	return &DeckService{
		db: database.GetDB(),
	}
}

// GetDB 获取数据库连接
func (s *DeckService) GetDB() *gorm.DB {
	return s.db
}

// CreateDeck 创建卡包
func (s *DeckService) CreateDeck(name string) (*models.Deck, error) {
	deck := &models.Deck{
		Name: name,
	}

	if err := s.db.Create(deck).Error; err != nil {
		return nil, err
	}

	return deck, nil
}

// GetDeckByID 根据ID获取卡包
func (s *DeckService) GetDeckByID(id uint) (*models.Deck, error) {
	var deck models.Deck
	if err := s.db.First(&deck, id).Error; err != nil {
		return nil, err
	}

	return &deck, nil
}

// GetAllDecks 获取所有卡包
func (s *DeckService) GetAllDecks() ([]models.Deck, error) {
	var decks []models.Deck
	if err := s.db.Find(&decks).Error; err != nil {
		return nil, err
	}

	return decks, nil
}

// UpdateDeck 更新卡包
func (s *DeckService) UpdateDeck(id uint, name string, archived *bool) (*models.Deck, error) {
	var deck models.Deck
	if err := s.db.First(&deck, id).Error; err != nil {
		return nil, err
	}

	if name != "" {
		deck.Name = name
	}

	if archived != nil {
		deck.Archived = *archived
	}

	if err := s.db.Save(&deck).Error; err != nil {
		return nil, err
	}

	return &deck, nil
}

// DeleteDeck 删除卡包
func (s *DeckService) DeleteDeck(id uint) error {
	if err := s.db.Delete(&models.Deck{}, id).Error; err != nil {
		return err
	}

	return nil
}

// GetDeckStats 获取卡包统计信息
func (s *DeckService) GetDeckStats(deckID uint) (*models.DeckStats, error) {
	stats := &models.DeckStats{}
	var count int64

	// 获取总卡片数
	if err := s.db.Model(&models.Card{}).Where("deck_id = ?", deckID).Count(&count).Error; err != nil {
		return nil, err
	}
	stats.TotalCards = int(count)

	// 获取待复习卡片数
	query := s.db.
		Table("cards").
		Select("COUNT(cards.id)").
		Joins("LEFT JOIN reviews ON cards.id = reviews.card_id").
		Where("cards.deck_id = ?", deckID).
		Where("(reviews.next_review <= date('now') OR reviews.id IS NULL)")

	if err := query.Scan(&count).Error; err != nil {
		return nil, err
	}
	stats.DueCards = int(count)

	// 获取标签数量
	if err := s.db.Model(&models.Tag{}).Where("deck_id = ?", deckID).Count(&count).Error; err != nil {
		return nil, err
	}
	stats.TagCount = int(count)

	// TODO: 实现今日和本周学习统计
	stats.TodayStudied = 0
	stats.WeekStudied = 0

	return stats, nil
}

// GetAllDecksWithStats 获取所有卡包及其统计信息
func (s *DeckService) GetAllDecksWithStats() ([]models.DeckWithStats, error) {
	var decks []models.Deck
	if err := s.db.Find(&decks).Error; err != nil {
		return nil, err
	}

	var result []models.DeckWithStats
	for _, deck := range decks {
		stats, err := s.GetDeckStats(deck.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, models.DeckWithStats{
			Deck:  deck,
			Stats: *stats,
		})
	}

	return result, nil
}
