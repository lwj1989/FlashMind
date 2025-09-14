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

// TestListTags 测试获取标签列表
func TestListTags(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)

	tag1 := models.Tag{DeckID: &deck.ID, Name: "测试标签1"}
	tag2 := models.Tag{DeckID: &deck.ID, Name: "测试标签2"}
	db.Create(&tag1)
	db.Create(&tag2)

	w := httptest.NewRecorder()
	// 使用创建的卡包ID
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/tags/deck/%d", deck.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "操作成功", response["message"])
	
	// 从data字段获取标签列表
	data := response["data"].([]interface{})
	assert.Equal(t, 2, len(data))
}

// TestCreateTag 测试创建标签
func TestCreateTag(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试卡包
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)
	
	// 确保卡包创建成功
	var createdDeck models.Deck
	db.First(&createdDeck, deck.ID)
	assert.Equal(t, "测试卡包", createdDeck.Name)

	tagData := map[string]interface{}{
		"name": "新标签",
	}
	jsonData, _ := json.Marshal(tagData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/tags/deck/%d", deck.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	// 打印响应内容以便调试
	t.Logf("Response status: %d", w.Code)
	t.Logf("Response body: %s", w.Body.String())
	
	// 先检查状态码
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		return
	}
	
	// 检查响应内容是否为有效的JSON
	if !json.Valid(w.Body.Bytes()) {
		t.Errorf("Response is not valid JSON: %s", w.Body.String())
		return
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "操作成功", response["message"])
	
	// 从data字段获取标签数据
	data := response["data"].(map[string]interface{})
	tag := data // 整个data就是标签对象
	assert.Equal(t, "新标签", tag["name"])
	assert.Equal(t, deck.ID, uint(tag["deck_id"].(float64)))
}

// TestGetTag 测试获取标签详情
func TestGetTag(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)

	tag := models.Tag{DeckID: &deck.ID, Name: "测试标签"}
	db.Create(&tag)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/tags/%d", tag.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "操作成功", response["message"])
	
	// 从data字段获取标签数据
	data := response["data"].(map[string]interface{})
	tagData := data // 整个data就是标签对象
	assert.Equal(t, "测试标签", tagData["name"])
}

// TestUpdateTag 测试更新标签
func TestUpdateTag(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)

	tag := models.Tag{DeckID: &deck.ID, Name: "测试标签"}
	db.Create(&tag)

	updateData := map[string]string{
		"name": "更新后的标签",
	}
	jsonData, _ := json.Marshal(updateData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/api/v1/tags/%d", tag.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "操作成功", response["message"])
	
	// 从data字段获取标签数据
	data := response["data"].(map[string]interface{})
	tagData := data // 整个data就是标签对象
	assert.Equal(t, "更新后的标签", tagData["name"])
}

// TestDeleteTag 测试删除标签
func TestDeleteTag(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)

	tag := models.Tag{DeckID: &deck.ID, Name: "测试标签"}
	db.Create(&tag)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/tags/%d", tag.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "操作成功", response["message"])

	// 验证标签已被删除
	var count int64
	db.Model(&models.Tag{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

// TestGetTagStats 测试获取标签统计信息
func TestGetTagStats(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)

	tag := models.Tag{DeckID: &deck.ID, Name: "测试标签"}
	db.Create(&tag)

	// 创建卡片
	card1 := models.Card{DeckID: deck.ID, TagID: &tag.ID, Question: "问题1", Answer: "答案1"}
	card2 := models.Card{DeckID: deck.ID, TagID: &tag.ID, Question: "问题2", Answer: "答案2"}
	db.Create(&card1)
	db.Create(&card2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/tags/%d/stats", tag.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "操作成功", response["message"])
	
	// 从data字段获取统计信息
	data := response["data"].(map[string]interface{})
	stats := data // 整个data就是统计信息对象
	assert.Equal(t, float64(2), stats["total_cards"])
}

// TestGetTagsByDeck 测试获取卡包下的标签
func TestGetTagsByDeck(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// 创建测试数据
	deck := models.Deck{Name: "测试卡包"}
	db.Create(&deck)

	tag1 := models.Tag{DeckID: &deck.ID, Name: "测试标签1"}
	tag2 := models.Tag{DeckID: &deck.ID, Name: "测试标签2"}
	db.Create(&tag1)
	db.Create(&tag2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/tags/deck/%d", deck.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, "SUCCESS", response["code"])
	assert.Equal(t, "操作成功", response["message"])
	
	// 从data字段获取标签列表
	data := response["data"].([]interface{})
	tags := data // 整个data就是标签列表
	assert.Equal(t, 2, len(tags))
}