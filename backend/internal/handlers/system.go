package handlers

import (
	"encoding/json"
	"flashcard/internal/models"
	"flashcard/internal/services"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// SystemHandler 系统管理处理器
type SystemHandler struct {
	deckService         *services.DeckService
	cardService         *services.CardService
	tagService          *services.TagService
	importExportService *services.ImportExportService
}

// NewSystemHandler 创建系统管理处理器实例
func NewSystemHandler() *SystemHandler {
	return &SystemHandler{
		deckService:         services.NewDeckService(),
		cardService:         services.NewCardService(),
		tagService:          services.NewTagService(),
		importExportService: services.NewImportExportService(),
	}
}

// BackupData 备份所有数据表
func (h *SystemHandler) BackupData(c *gin.Context) {
	// 创建临时目录
	tempDir := "./temp/backup_" + fmt.Sprintf("%d", time.Now().Unix())
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "创建临时目录失败", err.Error()))
		return
	}
	defer os.RemoveAll(tempDir)

	// 初始化备份数据结构
	backupData := models.BackupData{
		Version:    "1.0.0",
		ExportDate: time.Now(),
		Decks:      []models.DeckBackup{},
		Tags:       []models.TagBackup{},
		Cards:      []models.CardBackup{},
		Reviews:    []models.ReviewBackup{},
	}

	// 备份所有卡包
	var decks []models.Deck
	if err := h.deckService.GetDB().Find(&decks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "备份卡包失败", err.Error()))
		return
	}
	for _, deck := range decks {
		backupData.Decks = append(backupData.Decks, models.DeckBackup{
			ID:        deck.ID,
			Name:      deck.Name,
			Archived:  deck.Archived,
			CreatedAt: deck.CreatedAt,
			UpdatedAt: deck.UpdatedAt,
		})
	}

	// 备份所有标签
	var tags []models.Tag
	if err := h.tagService.GetDB().Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "备份标签失败", err.Error()))
		return
	}
	for _, tag := range tags {
		backupData.Tags = append(backupData.Tags, models.TagBackup{
			ID:        tag.ID,
			DeckID:    tag.DeckID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		})
	}

	// 备份所有卡片
	var cards []models.Card
	if err := h.cardService.GetDB().Find(&cards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "备份卡片失败", err.Error()))
		return
	}
	for _, card := range cards {
		backupData.Cards = append(backupData.Cards, models.CardBackup{
			ID:        card.ID,
			DeckID:    card.DeckID,
			TagID:     card.TagID,
			Question:  card.Question,
			Answer:    card.Answer,
			CreatedAt: card.CreatedAt,
			UpdatedAt: card.UpdatedAt,
		})
	}

	// 备份所有复习记录
	var reviews []models.Review
	if err := h.cardService.GetDB().Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "备份复习记录失败", err.Error()))
		return
	}
	for _, review := range reviews {
		backupData.Reviews = append(backupData.Reviews, models.ReviewBackup{
			ID:          review.ID,
			CardID:      review.CardID,
			EFactor:     review.EFactor,
			Interval:    review.Interval,
			Repetitions: review.Repetitions,
			NextReview:  review.NextReview,
			CreatedAt:   review.CreatedAt,
			UpdatedAt:   review.UpdatedAt,
		})
	}

	// 创建备份文件
	backupFilename := filepath.Join(tempDir, fmt.Sprintf("flashmind_complete_backup_%s.json", time.Now().Format("2006-01-02_15-04-05")))
	backupFile, err := os.Create(backupFilename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "创建备份文件失败", err.Error()))
		return
	}
	defer backupFile.Close()

	encoder := json.NewEncoder(backupFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(backupData); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "写入备份文件失败", err.Error()))
		return
	}

	// 直接下载文件
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(backupFilename))
	c.Header("Content-Type", "application/json")
	c.File(backupFilename)
}

// RestoreData 恢复所有数据表
func (h *SystemHandler) RestoreData(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "请选择要恢复的备份文件"))
		return
	}

	// 检查文件大小（限制为50MB）
	if file.Size > 50*1024*1024 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "文件大小不能超过50MB"))
		return
	}

	// 创建临时文件
	tempDir := "./temp"
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		os.Mkdir(tempDir, 0755)
	}

	tempFile, err := os.CreateTemp(tempDir, "restore-*"+filepath.Ext(file.Filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "创建临时文件失败", err.Error()))
		return
	}
	defer os.Remove(tempFile.Name())

	// 保存上传的文件
	if err := c.SaveUploadedFile(file, tempFile.Name()); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "保存文件失败", err.Error()))
		return
	}

	// 读取并解析备份文件
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "读取备份文件失败", err.Error()))
		return
	}

	var backupData models.BackupData
	if err := json.Unmarshal(content, &backupData); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的备份文件格式"))
		return
	}

	// 验证备份文件版本
	if backupData.Version == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "不支持的备份文件版本"))
		return
	}

	// 先清空现有数据（可选，根据需求）
	clearAll := c.DefaultQuery("clear_existing", "true")
	if clearAll == "true" {
		if err := h.clearAllDataInternal(); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "清空现有数据失败", err.Error()))
			return
		}
	}

	db := h.deckService.GetDB()
	restoredCounts := gin.H{
		"decks":   0,
		"tags":    0,
		"cards":   0,
		"reviews": 0,
	}

	// 开始数据库事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 恢复卡包数据
	for _, deckBackup := range backupData.Decks {
		deck := models.Deck{
			ID:        deckBackup.ID,
			Name:      deckBackup.Name,
			Archived:  deckBackup.Archived,
			CreatedAt: deckBackup.CreatedAt,
			UpdatedAt: deckBackup.UpdatedAt,
		}
		if err := tx.Create(&deck).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "恢复卡包失败", err.Error()))
			return
		}
		restoredCounts["decks"] = restoredCounts["decks"].(int) + 1
	}

	// 恢复标签数据
	for _, tagBackup := range backupData.Tags {
		tag := models.Tag{
			ID:        tagBackup.ID,
			DeckID:    tagBackup.DeckID,
			Name:      tagBackup.Name,
			CreatedAt: tagBackup.CreatedAt,
			UpdatedAt: tagBackup.UpdatedAt,
		}
		if err := tx.Create(&tag).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "恢复标签失败", err.Error()))
			return
		}
		restoredCounts["tags"] = restoredCounts["tags"].(int) + 1
	}

	// 恢复卡片数据
	for _, cardBackup := range backupData.Cards {
		card := models.Card{
			ID:        cardBackup.ID,
			DeckID:    cardBackup.DeckID,
			TagID:     cardBackup.TagID,
			Question:  cardBackup.Question,
			Answer:    cardBackup.Answer,
			CreatedAt: cardBackup.CreatedAt,
			UpdatedAt: cardBackup.UpdatedAt,
		}
		if err := tx.Create(&card).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "恢复卡片失败", err.Error()))
			return
		}
		restoredCounts["cards"] = restoredCounts["cards"].(int) + 1
	}

	// 恢复复习记录数据
	for _, reviewBackup := range backupData.Reviews {
		review := models.Review{
			ID:          reviewBackup.ID,
			CardID:      reviewBackup.CardID,
			EFactor:     reviewBackup.EFactor,
			Interval:    reviewBackup.Interval,
			Repetitions: reviewBackup.Repetitions,
			NextReview:  reviewBackup.NextReview,
			CreatedAt:   reviewBackup.CreatedAt,
			UpdatedAt:   reviewBackup.UpdatedAt,
		}
		if err := tx.Create(&review).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "恢复复习记录失败", err.Error()))
			return
		}
		restoredCounts["reviews"] = restoredCounts["reviews"].(int) + 1
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "恢复数据失败", err.Error()))
		return
	}

	response := models.SuccessResponse(restoredCounts)
	response.Message = fmt.Sprintf("成功恢复数据：%d个卡包，%d个标签，%d张卡片，%d条复习记录",
		restoredCounts["decks"], restoredCounts["tags"], restoredCounts["cards"], restoredCounts["reviews"])
	c.JSON(http.StatusOK, response)
}

// ClearAllData 清空所有数据
func (h *SystemHandler) ClearAllData(c *gin.Context) {
	if err := h.clearAllDataInternal(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "清空数据失败", err.Error()))
		return
	}

	response := models.SuccessResponse(gin.H{
		"message": "所有数据已清空",
	})
	response.Message = "数据清空成功"
	c.JSON(http.StatusOK, response)
}

// clearAllDataInternal 内部清空数据方法
func (h *SystemHandler) clearAllDataInternal() error {
	db := h.deckService.GetDB()

	// 开始事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 清空所有表的数据，按照外键依赖顺序
	// 1. 先删除复习记录
	if err := tx.Exec("DELETE FROM reviews").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("清空复习记录失败: %v", err)
	}

	// 2. 删除卡片
	if err := tx.Exec("DELETE FROM cards").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("清空卡片失败: %v", err)
	}

	// 3. 删除标签
	if err := tx.Exec("DELETE FROM tags").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("清空标签失败: %v", err)
	}

	// 4. 删除卡包
	if err := tx.Exec("DELETE FROM decks").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("清空卡包失败: %v", err)
	}

	// 重置自增ID（SQLite语法）
	tables := []string{"decks", "tags", "cards", "reviews"}
	for _, table := range tables {
		if err := tx.Exec(fmt.Sprintf("DELETE FROM sqlite_sequence WHERE name='%s'", table)).Error; err != nil {
			// 忽略错误，因为表可能没有自增字段
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("清空数据事务提交失败: %v", err)
	}

	return nil
}

// GetSystemStats 获取系统统计信息
func (h *SystemHandler) GetSystemStats(c *gin.Context) {
	// 获取统计信息
	decks, _ := h.deckService.GetAllDecks()

	// 获取总卡片数
	totalCards := int64(0)
	totalTags := 0
	for _, deck := range decks {
		cardsResponse, _ := h.cardService.GetCardsByDeckID(deck.ID, 1, 1000) // 假设每个卡包不超过1000张卡片
		if cardsResponse != nil {
			totalCards += cardsResponse.Total
		}

		tags, _ := h.tagService.GetTagsByDeckID(deck.ID)
		totalTags += len(tags)
	}

	stats := gin.H{
		"total_decks": len(decks),
		"total_cards": totalCards,
		"total_tags":  totalTags,
		"version":     "1.0.0",
	}

	c.JSON(http.StatusOK, models.SuccessResponse(stats))
}
