package services

import (
	"flashcard/internal/models"
	"flashcard/pkg/database"

	"gorm.io/gorm"
)

// TagService 标签服务
type TagService struct {
	db *gorm.DB
}

// NewTagService 创建标签服务实例
func NewTagService() *TagService {
	return &TagService{
		db: database.GetDB(),
	}
}

// GetDB 获取数据库连接
func (s *TagService) GetDB() *gorm.DB {
	return s.db
}

// CreateTag 创建标签
func (s *TagService) CreateTag(deckID *uint, name string) (*models.Tag, error) {
	tag := &models.Tag{
		DeckID: deckID,
		Name:   name,
	}

	if err := s.db.Create(tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}

// GetTagByID 根据ID获取标签
func (s *TagService) GetTagByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	if err := s.db.First(&tag, id).Error; err != nil {
		return nil, err
	}

	return &tag, nil
}

// GetTagsByDeckID 根据卡包ID获取标签列表
func (s *TagService) GetTagsByDeckID(deckID uint) ([]models.Tag, error) {
	var tags []models.Tag
	if err := s.db.Where("deck_id = ?", deckID).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

// UpdateTag 更新标签
func (s *TagService) UpdateTag(id uint, name string) (*models.Tag, error) {
	var tag models.Tag
	if err := s.db.First(&tag, id).Error; err != nil {
		return nil, err
	}

	tag.Name = name

	if err := s.db.Save(&tag).Error; err != nil {
		return nil, err
	}

	return &tag, nil
}

// UpdateTagWithDeck 更新标签（包含卡包ID）
func (s *TagService) UpdateTagWithDeck(id uint, name string, deckID *uint) (*models.Tag, error) {
	var tag models.Tag
	if err := s.db.First(&tag, id).Error; err != nil {
		return nil, err
	}

	tag.Name = name
	tag.DeckID = deckID

	if err := s.db.Save(&tag).Error; err != nil {
		return nil, err
	}

	return &tag, nil
}

// DeleteTag 删除标签
func (s *TagService) DeleteTag(id uint, deleteCards bool) error {
	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 如果需要同时删除卡片
	if deleteCards {
		// 先删除关联的复习记录
		if err := tx.Exec("DELETE FROM reviews WHERE card_id IN (SELECT id FROM cards WHERE tag_id = ?)", id).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 删除卡片
		if err := tx.Where("tag_id = ?", id).Delete(&models.Card{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// 仅将卡片的标签置为空
		if err := tx.Model(&models.Card{}).Where("tag_id = ?", id).Update("tag_id", nil).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 删除标签
	if err := tx.Delete(&models.Tag{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetTagStats 获取标签统计信息
func (s *TagService) GetTagStats(tagID uint) (*models.TagStats, error) {
	stats := &models.TagStats{}
	var count int64

	// 获取总卡片数
	if err := s.db.Model(&models.Card{}).Where("tag_id = ?", tagID).Count(&count).Error; err != nil {
		return nil, err
	}
	stats.TotalCards = int(count)

	// 获取待复习卡片数
	query := s.db.
		Table("cards").
		Select("COUNT(cards.id)").
		Joins("LEFT JOIN reviews ON cards.id = reviews.card_id").
		Where("cards.tag_id = ?", tagID).
		Where("(reviews.next_review <= date('now') OR reviews.id IS NULL)")

	if err := query.Scan(&count).Error; err != nil {
		return nil, err
	}
	stats.DueCards = int(count)

	// TODO: 实现今日学习统计
	stats.TodayStudied = 0

	return stats, nil
}

// GetTagsByDeckIDWithStats 根据卡包ID获取标签及其统计信息
func (s *TagService) GetTagsByDeckIDWithStats(deckID uint) ([]models.TagWithStats, error) {
	var tags []models.Tag
	if err := s.db.Where("deck_id = ?", deckID).Find(&tags).Error; err != nil {
		return nil, err
	}

	var result []models.TagWithStats
	for _, tag := range tags {
		stats, err := s.GetTagStats(tag.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, models.TagWithStats{
			Tag:   tag,
			Stats: *stats,
		})
	}

	return result, nil
}

// GetAllTags 获取所有标签
func (s *TagService) GetAllTags() ([]models.Tag, error) {
	var tags []models.Tag
	if err := s.db.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// GetAllTagsWithStats 获取所有标签及其统计信息
func (s *TagService) GetAllTagsWithStats() ([]models.TagWithStats, error) {
	var tags []models.Tag
	if err := s.db.Find(&tags).Error; err != nil {
		return nil, err
	}

	var result []models.TagWithStats
	for _, tag := range tags {
		stats, err := s.GetTagStats(tag.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, models.TagWithStats{
			Tag:   tag,
			Stats: *stats,
		})
	}

	return result, nil
}
