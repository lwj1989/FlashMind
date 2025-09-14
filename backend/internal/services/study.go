package services

import (
	"flashcard/internal/models"
	"flashcard/pkg/database"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"
)

// StudyService 学习服务
type StudyService struct {
	db *gorm.DB
}

// NewStudyService 创建学习服务实例
func NewStudyService() *StudyService {
	return &StudyService{
		db: database.GetDB(),
	}
}

// StartDeckStudy 开始学习卡包
func (s *StudyService) StartDeckStudy(deckID uint, limit int) (*models.StudySession, error) {
	var cards []models.Card
	err := s.db.Where("deck_id = ?", deckID).
		Preload("Deck").
		Preload("Tag").
		Preload("Review").
		Order("RANDOM()").
		Limit(limit).
		Find(&cards).Error
	if err != nil {
		return nil, err
	}

	return s.createStudySession(cards), nil
}

// StartTagStudy 开始学习标签
func (s *StudyService) StartTagStudy(tagID uint, limit int) (*models.StudySession, error) {
	var cards []models.Card
	err := s.db.Where("tag_id = ?", tagID).
		Preload("Deck").
		Preload("Tag").
		Preload("Review").
		Order("RANDOM()").
		Limit(limit).
		Find(&cards).Error
	if err != nil {
		return nil, err
	}

	return s.createStudySession(cards), nil
}

// StartRandomStudy 开始随机学习
func (s *StudyService) StartRandomStudy(limit int) (*models.StudySession, error) {
	var cards []models.Card
	err := s.db.Preload("Deck").
		Preload("Tag").
		Preload("Review").
		Order("RANDOM()").
		Limit(limit).
		Find(&cards).Error
	if err != nil {
		return nil, err
	}

	return s.createStudySession(cards), nil
}

// GetDueCards 获取到期复习的卡片
func (s *StudyService) GetDueCards(limit int) (*models.StudySession, error) {
	var cards []models.Card
	now := time.Now()

	// 获取到期的卡片（包括新卡片）
	err := s.db.Joins("LEFT JOIN reviews ON cards.id = reviews.card_id").
		Where("reviews.next_review <= ? OR reviews.id IS NULL", now).
		Preload("Deck").
		Preload("Tag").
		Preload("Review").
		Order("reviews.next_review ASC").
		Limit(limit).
		Find(&cards).Error
	if err != nil {
		return nil, err
	}

	return s.createStudySession(cards), nil
}

// createStudySession 创建学习会话
func (s *StudyService) createStudySession(cards []models.Card) *models.StudySession {
	queue := make([]models.StudyQueue, len(cards))

	for i, card := range cards {
		queue[i] = models.StudyQueue{
			CardID:   card.ID,
			Question: card.Question,
			Answer:   card.Answer,
			DeckName: card.Deck.Name,
		}

		if card.Tag != nil {
			queue[i].TagName = card.Tag.Name
		}
	}

	return &models.StudySession{
		Queue:     queue,
		Current:   0,
		Total:     len(queue),
		Completed: 0,
		StartTime: time.Now(),
	}
}

// SubmitReview 提交复习结果
func (s *StudyService) SubmitReview(cardID uint, result models.ReviewResult) (*models.ReviewResponse, error) {
	// 获取或创建复习记录
	var review models.Review
	err := s.db.Where("card_id = ?", cardID).First(&review).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建新的复习记录
			review = models.Review{
				CardID:      cardID,
				EFactor:     2.5,
				Interval:    0,
				Repetitions: 0,
				NextReview:  time.Now(),
			}
		} else {
			return nil, err
		}
	}

	// 使用SM-2算法更新复习参数
	s.updateReviewBySM2(&review, result)

	// 保存或更新复习记录
	if review.ID == 0 {
		err = s.db.Create(&review).Error
	} else {
		err = s.db.Save(&review).Error
	}
	if err != nil {
		return nil, err
	}

	// 构建响应
	response := &models.ReviewResponse{
		Success:    true,
		NextReview: review.NextReview,
		Interval:   review.Interval,
		Message:    s.getReviewMessage(result, review.Interval),
	}

	return response, nil
}

// updateReviewBySM2 使用SM-2算法更新复习参数
func (s *StudyService) updateReviewBySM2(review *models.Review, result models.ReviewResult) {
	switch result {
	case models.Again:
		// 忘记了，重新开始
		review.Repetitions = 0
		review.Interval = 0
		review.NextReview = time.Now().Add(10 * time.Minute) // 10分钟后重新复习
		review.EFactor = math.Max(1.3, review.EFactor-0.2)

	case models.Hard:
		// 模糊记得，稍微增加难度
		if review.Repetitions == 0 {
			review.Interval = 1
		} else if review.Repetitions == 1 {
			review.Interval = 3
		} else {
			review.Interval = int(float64(review.Interval) * review.EFactor)
		}
		review.Repetitions++
		review.NextReview = time.Now().AddDate(0, 0, review.Interval)
		review.EFactor = math.Max(1.3, review.EFactor-0.15)

	case models.Good:
		// 记得很好，正常间隔
		if review.Repetitions == 0 {
			review.Interval = 1
		} else if review.Repetitions == 1 {
			review.Interval = 6
		} else {
			review.Interval = int(float64(review.Interval) * review.EFactor)
		}
		review.Repetitions++
		review.NextReview = time.Now().AddDate(0, 0, review.Interval)
		// 稍微增加EFactor
		review.EFactor = math.Min(2.5, review.EFactor+0.1)
	}
}

// getReviewMessage 获取复习消息
func (s *StudyService) getReviewMessage(result models.ReviewResult, interval int) string {
	switch result {
	case models.Again:
		return "没关系，10分钟后我们再来复习这张卡片"
	case models.Hard:
		if interval == 1 {
			return "明天我们再来复习这张卡片"
		}
		return fmt.Sprintf("%d天后会再次复习这张卡片", interval)
	case models.Good:
		if interval == 1 {
			return "明天我们再来复习这张卡片"
		}
		return fmt.Sprintf("%d天后会再次复习这张卡片", interval)
	default:
		return "复习完成"
	}
}
