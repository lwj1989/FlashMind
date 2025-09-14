package services

import (
	"flashcard/internal/models"
	"flashcard/pkg/database"

	"gorm.io/gorm"
)

// CardService 卡片服务
type CardService struct {
	db *gorm.DB
}

// NewCardService 创建卡片服务实例
func NewCardService() *CardService {
	return &CardService{
		db: database.GetDB(),
	}
}

// GetDB 获取数据库连接
func (s *CardService) GetDB() *gorm.DB {
	return s.db
}

// CreateCard 创建卡片
func (s *CardService) CreateCard(deckID uint, tagID *uint, question, answer string) (*models.Card, error) {
	card := &models.Card{
		DeckID:   deckID,
		TagID:    tagID,
		Question: question,
		Answer:   answer,
	}

	if err := s.db.Create(card).Error; err != nil {
		return nil, err
	}

	return card, nil
}

// GetCardByID 根据ID获取卡片
func (s *CardService) GetCardByID(id uint) (*models.Card, error) {
	var card models.Card
	if err := s.db.Preload("Deck").Preload("Tag").First(&card, id).Error; err != nil {
		return nil, err
	}

	return &card, nil
}

// GetCardsByDeckID 根据卡包ID获取卡片列表
func (s *CardService) GetCardsByDeckID(deckID uint, page, pageSize int) (*models.CardListResponse, error) {
	var cards []models.Card
	var total int64

	// 获取总数
	if err := s.db.Model(&models.Card{}).Where("deck_id = ?", deckID).Count(&total).Error; err != nil {
		return nil, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := s.db.Where("deck_id = ?", deckID).
		Preload("Deck").
		Preload("Tag").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&cards).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	var cardResponses []models.CardResponse
	for _, card := range cards {
		cardResponse := models.CardResponse{
			Card: card,
		}
		if card.Tag != nil {
			cardResponse.TagName = card.Tag.Name
		}
		cardResponses = append(cardResponses, cardResponse)
	}

	// 计算总页数
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &models.CardListResponse{
		Cards:      cardResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetCardsByTagID 根据标签ID获取卡片列表
func (s *CardService) GetCardsByTagID(tagID uint, page, pageSize int) (*models.CardListResponse, error) {
	var cards []models.Card
	var total int64

	// 获取总数
	if err := s.db.Model(&models.Card{}).Where("tag_id = ?", tagID).Count(&total).Error; err != nil {
		return nil, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := s.db.Where("tag_id = ?", tagID).
		Preload("Deck").
		Preload("Tag").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&cards).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	var cardResponses []models.CardResponse
	for _, card := range cards {
		cardResponse := models.CardResponse{
			Card: card,
		}
		if card.Tag != nil {
			cardResponse.TagName = card.Tag.Name
		}
		cardResponses = append(cardResponses, cardResponse)
	}

	// 计算总页数
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &models.CardListResponse{
		Cards:      cardResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// SearchCards 搜索卡片
func (s *CardService) SearchCards(req models.CardSearchRequest) (*models.CardListResponse, error) {
	var cards []models.Card
	var total int64

	// 构建查询条件
	query := s.db.Model(&models.Card{})

	if req.DeckID != nil {
		query = query.Where("deck_id = ?", *req.DeckID)
	}

	if req.TagID != nil {
		query = query.Where("tag_id = ?", *req.TagID)
	}

	if req.Keyword != "" {
		query = query.Where("question LIKE ? OR answer LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 获取分页数据
	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("Deck").
		Preload("Tag").
		Offset(offset).
		Limit(req.PageSize).
		Order("created_at DESC").
		Find(&cards).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	var cardResponses []models.CardResponse
	for _, card := range cards {
		cardResponse := models.CardResponse{
			Card: card,
		}
		if card.Tag != nil {
			cardResponse.TagName = card.Tag.Name
		}
		cardResponses = append(cardResponses, cardResponse)
	}

	// 计算总页数
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &models.CardListResponse{
		Cards:      cardResponses,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateCard 更新卡片
func (s *CardService) UpdateCard(id uint, deckID uint, tagID *uint, question, answer string) (*models.Card, error) {
	var card models.Card
	if err := s.db.First(&card, id).Error; err != nil {
		return nil, err
	}

	card.DeckID = deckID
	card.TagID = tagID
	card.Question = question
	card.Answer = answer

	if err := s.db.Save(&card).Error; err != nil {
		return nil, err
	}

	return &card, nil
}

// DeleteCard 删除卡片
func (s *CardService) DeleteCard(id uint) error {
	if err := s.db.Delete(&models.Card{}, id).Error; err != nil {
		return err
	}

	return nil
}
