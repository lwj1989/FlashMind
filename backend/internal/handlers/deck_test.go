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

// TestListDecks 测试获取卡包列表
func TestListDecks(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck1 := models.Deck{Name: "测试卡包1"}
	deck2 := models.Deck{Name: "测试卡包2"}
	db.Create(&deck1)
	db.Create(&deck2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/decks", nil)
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
	
	// 获取decks字段
	decks, ok := data["decks"].([]interface{})
	assert.True(t, ok, "Decks should be an array")
	assert.Equal(t, 2, len(decks))
}

// TestCreateDeck 测试创建卡包
func TestCreateDeck(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	deckData := map[string]string{
		"name": "新卡包",
	}
	jsonData, _ := json.Marshal(deckData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/decks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "操作成功", response["message"])
	
	// 获取data字段，data就是deck对象
	deck, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "Response data should be a map")
	assert.Equal(t, "新卡包", deck["name"])
}

// TestGetDeck 测试获取卡包详情
func TestGetDeck(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/decks/%d", deck.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "操作成功", response["message"])
	
	// 获取data字段，data就是deck对象
	deckData, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "Response data should be a map")
	assert.Equal(t, "测试卡包", deckData["name"])
}

// TestUpdateDeck 测试更新卡包
func TestUpdateDeck(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)

	updateData := map[string]string{
		"name": "更新后的卡包",
	}
	jsonData, _ := json.Marshal(updateData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/api/v1/decks/%d", deck.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "操作成功", response["message"])
	
	// 获取data字段，data就是deck对象
	deckData, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "Response data should be a map")
	assert.Equal(t, "更新后的卡包", deckData["name"])
}

// TestDeleteDeck 测试删除卡包
func TestDeleteDeck(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/decks/%d", deck.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 验证卡包已被删除
	var count int64
	db.Model(&models.Deck{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

// TestGetDeckStats 测试获取卡包统计信息
func TestGetDeckStats(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)

	// 创建标签
	tag := models.Tag{DeckID: &deck.ID, Name: "测试标签"}
	db.Create(&tag)

	// 创建卡片
	card1 := models.Card{DeckID: deck.ID, TagID: &tag.ID, Question: "问题1", Answer: "答案1"}
	card2 := models.Card{DeckID: deck.ID, Question: "问题2", Answer: "答案2"}
	db.Create(&card1)
	db.Create(&card2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/decks/%d/stats", deck.ID), nil)
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
	
	// 获取stats字段
	stats, ok := data["stats"].(map[string]interface{})
	assert.True(t, ok, "Stats should be a map")
	assert.Equal(t, float64(2), stats["total_cards"])
	assert.Equal(t, float64(1), stats["tag_count"])
}