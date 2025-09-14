package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"flashcard/internal/models"
)

// TestSearchCards 测试搜索卡片
func TestSearchCards(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)

	tag := models.Tag{DeckID: &deck.ID, Name: "测试标签"}
	db.Create(&tag)

	card1 := models.Card{DeckID: deck.ID, TagID: &tag.ID, Question: "Go语言是什么", Answer: "Go是一种编程语言"}
	card2 := models.Card{DeckID: deck.ID, Question: "Python是什么", Answer: "Python是一种编程语言"}
	db.Create(&card1)
	db.Create(&card2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cards?keyword=Go&page=1&page_size=20", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "操作成功", response["message"])
	
	// 获取data字段
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "Response data should be a map")
	
	// 获取cards字段
	cards, ok := data["cards"].([]interface{})
	assert.True(t, ok, "Cards should be an array")
	
	// 验证卡片数量
	assert.Equal(t, 1, len(cards))

	card := cards[0].(map[string]interface{})
	assert.Equal(t, "Go语言是什么", card["question"])
}

// TestCreateCard 测试创建卡片
func TestCreateCard(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)

	tag := models.Tag{DeckID: &deck.ID, Name: "测试标签"}
	db.Create(&tag)

	cardData := map[string]interface{}{
		"deck_id":  deck.ID,
		"tag_id":   tag.ID,
		"question": "什么是Go语言？",
		"answer":   "Go是一种开源的编程语言",
	}
	jsonData, _ := json.Marshal(cardData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/cards", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "操作成功", response["message"])
	
	// 获取data字段，data就是card对象
	card, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "Response data should be a map")
	
	// 验证卡片内容
	assert.Equal(t, "什么是Go语言？", card["question"])
	assert.Equal(t, "Go是一种开源的编程语言", card["answer"])
	assert.Equal(t, deck.ID, uint(card["deck_id"].(float64)))
	assert.Equal(t, tag.ID, uint(card["tag_id"].(float64)))
}

// TestGetCard 测试获取卡片详情
func TestGetCard(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)

	card := models.Card{DeckID: deck.ID, Question: "测试问题", Answer: "测试答案"}
	db.Create(&card)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/cards/%d", card.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "操作成功", response["message"])
	
	// 获取data字段，data就是card对象
	cardData, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "Response data should be a map")
	
	// 验证卡片内容
	assert.Equal(t, "测试问题", cardData["question"])
	assert.Equal(t, "测试答案", cardData["answer"])
}

// TestUpdateCard 测试更新卡片
func TestUpdateCard(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)

	card := models.Card{DeckID: deck.ID, Question: "测试问题", Answer: "测试答案"}
	db.Create(&card)

	updateData := map[string]interface{}{
		"deck_id":  deck.ID,
		"question": "更新后的问题",
		"answer":   "更新后的答案",
	}
	jsonData, _ := json.Marshal(updateData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/api/v1/cards/%d", card.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "操作成功", response["message"])
	
	// 获取data字段，data就是card对象
	cardData, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "Response data should be a map")
	
	// 验证卡片内容
	assert.Equal(t, "更新后的问题", cardData["question"])
	assert.Equal(t, "更新后的答案", cardData["answer"])
}

// TestDeleteCard 测试删除卡片
func TestDeleteCard(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)

	card := models.Card{DeckID: deck.ID, Question: "测试问题", Answer: "测试答案"}
	db.Create(&card)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/cards/%d", card.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 验证卡片已被删除
	var count int64
	db.Model(&models.Card{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

// TestGetCardsByDeck 测试获取卡包下的卡片
func TestGetCardsByDeck(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)

	card1 := models.Card{DeckID: deck.ID, Question: "问题1", Answer: "答案1"}
	card2 := models.Card{DeckID: deck.ID, Question: "问题2", Answer: "答案2"}
	db.Create(&card1)
	db.Create(&card2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/cards?deck_id=%d&page=1&page_size=20", deck.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "操作成功", response["message"])
	
	// 获取data字段
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "Response data should be a map")
	
	// 获取cards字段
	cards, ok := data["cards"].([]interface{})
	assert.True(t, ok, "Cards should be an array")
	
	// 验证卡片数量
	assert.Equal(t, 2, len(cards))
}

// TestGetCardsByTag 测试获取标签下的卡片
func TestGetCardsByTag(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)

	tag := models.Tag{DeckID: &deck.ID, Name: "测试标签"}
	db.Create(&tag)

	card1 := models.Card{DeckID: deck.ID, TagID: &tag.ID, Question: "问题1", Answer: "答案1"}
	card2 := models.Card{DeckID: deck.ID, TagID: &tag.ID, Question: "问题2", Answer: "答案2"}
	db.Create(&card1)
	db.Create(&card2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/tags/%d/cards?page=1&page_size=20", tag.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "操作成功", response["message"])
	
	// 获取data字段
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "Response data should be a map")
	
	// 获取cards字段
	cards, ok := data["cards"].([]interface{})
	assert.True(t, ok, "Cards should be an array")
	
	// 验证卡片数量
	assert.Equal(t, 2, len(cards))
}