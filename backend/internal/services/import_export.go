package services

import (
	"encoding/csv"
	"encoding/json"
	"flashcard/internal/models"
	"flashcard/pkg/database"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ImportExportService 导入导出服务
type ImportExportService struct {
	db *gorm.DB
}

// NewImportExportService 创建导入导出服务实例
func NewImportExportService() *ImportExportService {
	return &ImportExportService{
		db: database.GetDB(),
	}
}

// GetDB 获取数据库连接
func (s *ImportExportService) GetDB() *gorm.DB {
	return s.db
}

// ExportDeck 导出卡包为JSON格式
func (s *ImportExportService) ExportDeck(deckID uint, format string) (string, error) {
	// 获取卡包信息
	var deck models.Deck
	if err := s.db.First(&deck, deckID).Error; err != nil {
		return "", err
	}

	// 获取卡包下的所有标签
	var tags []models.Tag
	if err := s.db.Where("deck_id = ?", deckID).Find(&tags).Error; err != nil {
		return "", err
	}

	// 获取卡包下的所有卡片
	var cards []models.Card
	if err := s.db.Where("deck_id = ?", deckID).Preload("Tag").Find(&cards).Error; err != nil {
		return "", err
	}

	// 构建导出数据结构（按照导入时的格式）
	var tagExports []models.TagExport
	for _, tag := range tags {
		tagExports = append(tagExports, models.TagExport{
			Name:        tag.Name,
			Description: "", // 目前模型中没有描述字段
			CreatedAt:   tag.CreatedAt,
		})
	}

	var cardExports []models.CardExport
	for _, card := range cards {
		cardExport := models.CardExport{
			Question:  card.Question,
			Answer:    card.Answer,
			CreatedAt: card.CreatedAt,
		}
		if card.Tag != nil {
			cardExport.TagName = card.Tag.Name
		}
		cardExports = append(cardExports, cardExport)
	}

	exportData := models.DeckExport{
		Name:        deck.Name,
		Description: "", // 目前模型中没有描述字段
		Cards:       cardExports,
		Tags:        tagExports,
		CreatedAt:   deck.CreatedAt,
	}

	// 确保导出目录存在
	exportDir := "./data"
	if _, err := os.Stat(exportDir); os.IsNotExist(err) {
		os.Mkdir(exportDir, 0755)
	}

	// 根据格式导出
	switch strings.ToLower(format) {
	case "json":
		return s.exportToJSON(exportData, deck.Name)
	case "csv":
		return s.exportToCSV(exportData, deck.Name)
	case "txt":
		return s.exportToTXT(exportData, deck.Name)
	default:
		return "", fmt.Errorf("不支持的导出格式: %s", format)
	}
}

// exportToJSON 导出为JSON格式
func (s *ImportExportService) exportToJSON(data interface{}, deckName string) (string, error) {
	// 创建文件名
	filename := fmt.Sprintf("./data/%s_%s.json", deckName, getCurrentTimestamp())
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 编码JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return "", err
	}

	return filename, nil
}

// exportToCSV 导出为CSV格式
func (s *ImportExportService) exportToCSV(data interface{}, deckName string) (string, error) {
	// 创建文件名
	filename := fmt.Sprintf("./data/%s_%s.csv", deckName, getCurrentTimestamp())
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建CSV写入器
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	if err := writer.Write([]string{"ID", "Question", "Answer", "Tag"}); err != nil {
		return "", err
	}

	// 转换数据并写入CSV
	exportData, ok := data.(models.DeckExport)
	if !ok {
		return "", fmt.Errorf("数据格式错误")
	}

	// 写入卡片数据
	for i, card := range exportData.Cards {
		if err := writer.Write([]string{
			strconv.Itoa(i + 1), // 使用序号作为ID
			card.Question,
			card.Answer,
			card.TagName,
		}); err != nil {
			return "", err
		}
	}

	return filename, nil
}

// ImportDeck 从文件导入卡包
func (s *ImportExportService) ImportDeck(filePath string) (*models.Deck, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("文件不存在: %s", filePath)
	}

	// 根据文件扩展名确定导入格式
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".json":
		return s.importFromJSON(filePath)
	case ".csv":
		return s.importFromCSV(filePath)
	case ".txt":
		return s.importFromTXT(filePath)
	default:
		return nil, fmt.Errorf("不支持的导入格式: %s", ext)
	}
}

// importFromJSON 从JSON文件导入
func (s *ImportExportService) importFromJSON(filePath string) (*models.Deck, error) {
	// 读取文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 尝试解析为 DeckExport 格式
	var deckExport models.DeckExport
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&deckExport); err != nil {
		// 如果解析失败，尝试旧格式
		file.Seek(0, 0)
		var importData struct {
			Deck  models.Deck   `json:"deck"`
			Tags  []models.Tag  `json:"tags"`
			Cards []models.Card `json:"cards"`
		}
		if err := json.NewDecoder(file).Decode(&importData); err != nil {
			return nil, fmt.Errorf("无法解析导入文件格式: %v", err)
		}
		// 转换为 DeckExport 格式
		deckExport = models.DeckExport{
			Name:        importData.Deck.Name,
			Description: "",
			CreatedAt:   importData.Deck.CreatedAt,
		}
		for _, tag := range importData.Tags {
			deckExport.Tags = append(deckExport.Tags, models.TagExport{
				Name:        tag.Name,
				Description: "",
				CreatedAt:   tag.CreatedAt,
			})
		}
		for _, card := range importData.Cards {
			cardExport := models.CardExport{
				Question:  card.Question,
				Answer:    card.Answer,
				CreatedAt: card.CreatedAt,
			}
			if card.TagID != nil {
				// 查找标签名称
				for _, tag := range importData.Tags {
					if tag.ID == *card.TagID {
						cardExport.TagName = tag.Name
						break
					}
				}
			}
			deckExport.Cards = append(deckExport.Cards, cardExport)
		}
	}

	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建卡包
	deck := models.Deck{
		Name:     deckExport.Name,
		Archived: false,
	}
	if err := tx.Create(&deck).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建标签映射（标签名到ID）
	tagMap := make(map[string]uint)
	for _, tagExport := range deckExport.Tags {
		newTag := models.Tag{
			DeckID: &deck.ID,
			Name:   tagExport.Name,
		}
		if err := tx.Create(&newTag).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		tagMap[tagExport.Name] = newTag.ID
	}

	// 创建卡片
	cardCount := 0
	for _, cardExport := range deckExport.Cards {
		newCard := models.Card{
			DeckID:   deck.ID,
			Question: cardExport.Question,
			Answer:   cardExport.Answer,
		}
		// 如果卡片有关联标签，则使用新标签ID
		if cardExport.TagName != "" {
			if newTagID, exists := tagMap[cardExport.TagName]; exists {
				newCard.TagID = &newTagID
			}
		}
		if err := tx.Create(&newCard).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		cardCount++
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &deck, nil
}

// importFromCSV 从CSV文件导入
func (s *ImportExportService) importFromCSV(filePath string) (*models.Deck, error) {
	// 读取文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 创建CSV读取器
	reader := csv.NewReader(file)

	// 读取表头
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// 验证表头
	if len(header) < 3 || header[0] != "ID" || header[1] != "Question" || header[2] != "Answer" {
		return nil, fmt.Errorf("CSV格式错误，表头应为: ID,Question,Answer,Tag")
	}

	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建卡包（使用文件名作为卡包名）
	baseName := filepath.Base(filePath)
	deckName := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	deckName = strings.TrimSuffix(deckName, "_"+getCurrentTimestamp()) // 移除时间戳（如果有）

	deck := models.Deck{
		Name: deckName,
	}
	if err := tx.Create(&deck).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建标签映射（标签名到ID）
	tagMap := make(map[string]uint)

	// 读取数据行
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// 跳过空行
		if len(record) < 3 {
			continue
		}

		// 创建卡片
		card := models.Card{
			DeckID:   deck.ID,
			Question: record[1],
			Answer:   record[2],
		}

		// 如果有标签列，处理标签
		if len(record) > 3 && record[3] != "" {
			tagName := record[3]
			var tagID uint
			var exists bool

			// 检查标签是否已存在
			if tagID, exists = tagMap[tagName]; !exists {
				// 创建新标签
				tag := models.Tag{
					DeckID: &deck.ID,
					Name:   tagName,
				}
				if err := tx.Create(&tag).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
				tagID = tag.ID
				tagMap[tagName] = tagID
			}

			card.TagID = &tagID
		}

		// 创建卡片
		if err := tx.Create(&card).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &deck, nil
}

// exportToTXT 导出为TXT格式
func (s *ImportExportService) exportToTXT(data interface{}, deckName string) (string, error) {
	// 创建文件名
	filename := fmt.Sprintf("./data/%s_%s.txt", deckName, getCurrentTimestamp())
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 转换数据
	exportData, ok := data.(models.DeckExport)
	if !ok {
		return "", fmt.Errorf("数据格式错误")
	}

	// 按标签分组卡片
	tagCards := make(map[string][]models.CardExport)
	noTagCards := []models.CardExport{}

	for _, card := range exportData.Cards {
		if card.TagName != "" {
			tagCards[card.TagName] = append(tagCards[card.TagName], card)
		} else {
			noTagCards = append(noTagCards, card)
		}
	}

	// 写入有标签的卡片
	for tagName, cards := range tagCards {
		_, err := file.WriteString(fmt.Sprintf("# %s\n", tagName))
		if err != nil {
			return "", err
		}

		for i, card := range cards {
			_, err := file.WriteString(fmt.Sprintf("%s\n---\n%s", card.Question, card.Answer))
			if err != nil {
				return "", err
			}

			// 如果不是最后一张卡片，添加分隔符
			if i < len(cards)-1 {
				_, err := file.WriteString("\n===\n")
				if err != nil {
					return "", err
				}
			}
		}

		// 标签组之间添加换行
		_, err = file.WriteString("\n\n")
		if err != nil {
			return "", err
		}
	}

	// 写入无标签的卡片
	if len(noTagCards) > 0 {
		for i, card := range noTagCards {
			_, err := file.WriteString(fmt.Sprintf("%s\n---\n%s", card.Question, card.Answer))
			if err != nil {
				return "", err
			}

			// 如果不是最后一张卡片，添加分隔符
			if i < len(noTagCards)-1 {
				_, err := file.WriteString("\n===\n")
				if err != nil {
					return "", err
				}
			}
		}
	}

	return filename, nil
}

// getCurrentTimestamp 获取当前时间戳（年月日格式）
func getCurrentTimestamp() string {
	return time.Now().Format("20060102")
}

// importFromTXT 从TXT文件导入
func (s *ImportExportService) importFromTXT(filePath string) (*models.Deck, error) {
	// 读取文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建卡包（使用文件名作为卡包名）
	baseName := filepath.Base(filePath)
	deckName := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	deckName = strings.TrimSuffix(deckName, "_"+getCurrentTimestamp()) // 移除时间戳（如果有）

	deck := models.Deck{
		Name: deckName,
	}
	if err := tx.Create(&deck).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 解析TXT内容
	contentStr := string(content)

	// 创建标签映射（标签名到ID）
	tagMap := make(map[string]uint)
	var currentTagID *uint

	// 按===分割卡片
	cardParts := strings.Split(contentStr, "===")

	for _, cardPart := range cardParts {
		cardPart = strings.TrimSpace(cardPart)
		if cardPart == "" {
			continue
		}

		// 检查每个卡片部分是否有标签
		lines := strings.Split(cardPart, "\n")
		var cleanedLines []string

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// 检查是否是标签行
			if strings.HasPrefix(line, "#") {
				tagName := strings.TrimSpace(line[1:])
				if tagName != "" {
					// 检查标签是否已存在
					var exists bool
					var tagID uint

					if tagID, exists = tagMap[tagName]; !exists {
						// 创建新标签
						tag := models.Tag{
							DeckID: &deck.ID,
							Name:   tagName,
						}
						if err := tx.Create(&tag).Error; err != nil {
							tx.Rollback()
							return nil, err
						}
						tagID = tag.ID
						tagMap[tagName] = tagID
					}
					currentTagID = &tagID
				}
				// 跳过标签行，不加入到cleanedLines中
				continue
			}

			cleanedLines = append(cleanedLines, line)
		}

		// 重新组装卡片内容
		cleanedContent := strings.Join(cleanedLines, "\n")

		// 分割问题和答案
		parts := strings.Split(cleanedContent, "---")
		if len(parts) != 2 {
			continue // 跳过格式不正确的卡片
		}

		questionText := strings.TrimSpace(parts[0])
		answerText := strings.TrimSpace(parts[1])

		// 跳过空的问题或答案
		if questionText == "" || answerText == "" {
			continue
		}

		// 创建卡片
		card := models.Card{
			DeckID:   deck.ID,
			Question: questionText,
			Answer:   answerText,
			TagID:    currentTagID,
		}

		// 创建卡片
		if err := tx.Create(&card).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &deck, nil
}
