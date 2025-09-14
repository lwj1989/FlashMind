package handlers

import (
	"flashcard/internal/models"
	"flashcard/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CardHandler 卡片处理器
type CardHandler struct {
	cardService *services.CardService
}

// NewCardHandler 创建卡片处理器实例
func NewCardHandler() *CardHandler {
	return &CardHandler{
		cardService: services.NewCardService(),
	}
}

// CreateCard 创建卡片
func (h *CardHandler) CreateCard(c *gin.Context) {
	var req struct {
		DeckID   uint   `json:"deck_id" binding:"required"`
		TagID    *uint  `json:"tag_id"`
		Question string `json:"question" binding:"required"`
		Answer   string `json:"answer" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "请求参数格式错误"))
		return
	}

	card, err := h.cardService.CreateCard(req.DeckID, req.TagID, req.Question, req.Answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "创建卡片失败", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse(card))
}

// GetCard 获取单个卡片
func (h *CardHandler) GetCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的卡片ID"))
		return
	}

	card, err := h.cardService.GetCardByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(models.CodeNotFound, "卡片不存在"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(card))
}

// GetCardsByDeck 获取卡包下的所有卡片
func (h *CardHandler) GetCardsByDeck(c *gin.Context) {
	deckIDStr := c.Param("id")
	deckID, err := strconv.ParseUint(deckIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的卡包ID"))
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 验证分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	cards, err := h.cardService.GetCardsByDeckID(uint(deckID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "获取卡片列表失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(cards))
}

// GetCardsByTag 获取标签下的所有卡片
func (h *CardHandler) GetCardsByTag(c *gin.Context) {
	tagIDStr := c.Param("id")
	tagID, err := strconv.ParseUint(tagIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的标签ID"))
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 验证分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	cards, err := h.cardService.GetCardsByTagID(uint(tagID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "获取卡片列表失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(cards))
}

// SearchCards 搜索卡片
func (h *CardHandler) SearchCards(c *gin.Context) {
	var req models.CardSearchRequest

	// 绑定查询参数
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "请求参数格式错误"))
		return
	}

	// 验证分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	cards, err := h.cardService.SearchCards(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "搜索卡片失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(cards))
}

// UpdateCard 更新卡片
func (h *CardHandler) UpdateCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的卡片ID"))
		return
	}

	var req struct {
		DeckID   uint   `json:"deck_id" binding:"required"`
		TagID    *uint  `json:"tag_id"`
		Question string `json:"question" binding:"required"`
		Answer   string `json:"answer" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "请求参数格式错误"))
		return
	}

	card, err := h.cardService.UpdateCard(uint(id), req.DeckID, req.TagID, req.Question, req.Answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "更新卡片失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(card))
}

// DeleteCard 删除卡片
func (h *CardHandler) DeleteCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的卡片ID"))
		return
	}

	if err := h.cardService.DeleteCard(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "删除卡片失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(nil))
}