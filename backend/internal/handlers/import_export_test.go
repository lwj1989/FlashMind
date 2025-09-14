package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"flashcard/internal/models"
)
func TestImportDeck(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试JSON数据
	deckData := map[string]interface{}{
		"deck": map[string]interface{}{
			"name":     "导入测试卡包",
			"archived": false,
		},
		"tags": []map[string]interface{}{},
		"cards": []map[string]interface{}{
			{
				"question": "问题1",
				"answer":   "答案1",
			},
			{
				"question": "问题2",
				"answer":   "答案2",
			},
		},
	}
	jsonData, _ := json.Marshal(deckData)

	// 创建multipart表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.json")
	part.Write(jsonData)
	writer.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/import-export/decks", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "导入成功", response["message"])

	// 验证数据已导入
	var deck models.Deck
	db.First(&deck, "name = ?", "导入测试卡包")
	assert.Equal(t, "导入测试卡包", deck.Name)

	var count int64
	db.Model(&models.Card{}).Where("deck_id = ?", deck.ID).Count(&count)
	assert.Equal(t, int64(2), count)
}

// TestExportDeck 测试导出卡包
func TestExportDeck(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "导出测试卡包"}
	db.Create(&deck)

	card1 := models.Card{DeckID: deck.ID, Question: "问题1", Answer: "答案1"}
	card2 := models.Card{DeckID: deck.ID, Question: "问题2", Answer: "答案2"}
	db.Create(&card1)
	db.Create(&card2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/import-export/decks/%d?format=json", deck.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	// 验证响应结构
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "卡包导出成功", response["message"])
	
	// 获取data字段
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "Response data should be a map")
	
	// 验证filename字段
	filename, ok := data["filename"].(string)
	assert.True(t, ok, "Filename should be a string")
	assert.True(t, len(filename) > 0, "Filename should not be empty")
}

// TestExportDeckCSV 测试导出卡包为CSV格式
func TestExportDeckCSV(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "导出测试卡包"}
	db.Create(&deck)

	card1 := models.Card{DeckID: deck.ID, Question: "问题1", Answer: "答案1"}
	card2 := models.Card{DeckID: deck.ID, Question: "问题2", Answer: "答案2"}
	db.Create(&card1)
	db.Create(&card2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/import-export/decks/%d?format=csv", deck.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	// 验证响应结构
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "卡包导出成功", response["message"])
	
	// 获取data字段
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "Response data should be a map")
	
	// 验证filename字段
	filename, ok := data["filename"].(string)
	assert.True(t, ok, "Filename should be a string")
	assert.True(t, len(filename) > 0, "Filename should not be empty")
	assert.True(t, strings.HasSuffix(filename, ".csv"), "Filename should have .csv extension")
}

// TestImportDeckCSV 测试导入CSV格式的卡包
func TestImportDeckCSV(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试CSV数据
	csvData := "ID,Question,Answer\n1,问题1,答案1\n2,问题2,答案2\n"

	// 创建multipart表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.csv")
	part.Write([]byte(csvData))
	writer.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/import-export/decks", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "导入成功", response["message"])

	// 验证数据已导入
	var deck models.Deck
	result := db.First(&deck, "name LIKE ?", "%upload%")
	if result.Error != nil {
		t.Logf("查找卡包失败: %v", result.Error)
		// 打印所有卡包
		var allDecks []models.Deck
		db.Find(&allDecks)
		t.Logf("数据库中的所有卡包:")
		for _, d := range allDecks {
			t.Logf("ID: %d, Name: %s", d.ID, d.Name)
		}
	}
	assert.True(t, strings.Contains(deck.Name, "upload"))

	var count int64
	db.Model(&models.Card{}).Where("deck_id = ?", deck.ID).Count(&count)
	assert.Equal(t, int64(2), count)
}

// TestImportDeckInvalidFormat 测试导入无效格式的文件
func TestImportDeckInvalidFormat(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建无效的测试数据
	invalidData := "这不是有效的JSON或CSV格式"

	// 创建multipart表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.xml") // 使用不支持的XML格式
	part.Write([]byte(invalidData))
	writer.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/import-export/decks", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "INVALID_PARAM", response["code"])
	assert.Contains(t, response["message"], "导入失败")
}

// TestExportDeckNotFound 测试导出不存在的卡包
func TestExportDeckNotFound(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/import-export/decks/999?format=json", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "NOT_FOUND", response["code"])
	assert.Contains(t, response["message"], "导出卡包失败")
}

// TestImportDeckTXT 测试导入TXT格式的卡包
func TestImportDeckTXT(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试TXT数据
	txtData := `问题1===答案1
问题2===答案2
# 标签1
问题3===答案3
问题4===答案4
# 标签2
问题5===答案5`

	// 创建multipart表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte(txtData))
	writer.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/import-export/decks", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "导入成功", response["message"])

	// 验证数据已导入
	var deck models.Deck
	result := db.First(&deck, "name LIKE ?", "%upload%")
	if result.Error != nil {
		t.Logf("查找卡包失败: %v", result.Error)
		// 打印所有卡包
		var allDecks []models.Deck
		db.Find(&allDecks)
		t.Logf("数据库中的所有卡包:")
		for _, d := range allDecks {
			t.Logf("ID: %d, Name: %s", d.ID, d.Name)
		}
	}
	assert.True(t, strings.Contains(deck.Name, "upload"))

	// 验证卡片数量
	var count int64
	db.Model(&models.Card{}).Where("deck_id = ?", deck.ID).Count(&count)
	assert.Equal(t, int64(5), count)

	// 验证标签数量
	var tagCount int64
	db.Model(&models.Tag{}).Where("deck_id = ?", deck.ID).Count(&tagCount)
	assert.Equal(t, int64(2), tagCount)

	// 验证标签名称
	var tags []models.Tag
	db.Where("deck_id = ?", deck.ID).Find(&tags)
	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}
	assert.Contains(t, tagNames, "标签1")
	assert.Contains(t, tagNames, "标签2")
}