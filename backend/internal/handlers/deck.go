package handlers

import (
	"flashcard/internal/models"
	"flashcard/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeckHandler 卡包处理器
type DeckHandler struct {
	deckService *services.DeckService
}

// NewDeckHandler 创建卡包处理器实例
func NewDeckHandler() *DeckHandler {
	return &DeckHandler{
		deckService: services.NewDeckService(),
	}
}

// CreateDeck 创建卡包
func (h *DeckHandler) CreateDeck(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required,min=1,max=100"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "卡包名称不能为空且长度应在1-100之间"))
		return
	}

	deck, err := h.deckService.CreateDeck(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "创建卡包失败", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse(deck))
}

// GetDecks 获取所有卡包
func (h *DeckHandler) GetDecks(c *gin.Context) {
	includeStats := c.DefaultQuery("include_stats", "false")

	if includeStats == "true" {
		decks, err := h.deckService.GetAllDecksWithStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "获取卡包列表失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse(map[string]interface{}{
			"decks": decks,
		}))
	} else {
		decks, err := h.deckService.GetAllDecks()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "获取卡包列表失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse(map[string]interface{}{
			"decks": decks,
		}))
	}
}

// GetDeck 获取单个卡包
func (h *DeckHandler) GetDeck(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的卡包ID"))
		return
	}

	deck, err := h.deckService.GetDeckByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(models.CodeNotFound, "卡包不存在"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(deck))
}

// UpdateDeck 更新卡包
func (h *DeckHandler) UpdateDeck(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的卡包ID"))
		return
	}

	var req struct {
		Name     *string `json:"name"`
		Archived *bool   `json:"archived"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "请求参数格式错误"))
		return
	}

	// 检查至少有一个字段需要更新
	if req.Name == nil && req.Archived == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "至少需要提供一个更新字段"))
		return
	}

	var name string
	if req.Name != nil {
		name = *req.Name
	}

	deck, err := h.deckService.UpdateDeck(uint(id), name, req.Archived)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "更新卡包失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(deck))
}

// DeleteDeck 删除卡包
func (h *DeckHandler) DeleteDeck(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的卡包ID"))
		return
	}

	if err := h.deckService.DeleteDeck(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "删除卡包失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(nil))
}

// GetDeckStats 获取卡包统计信息
func (h *DeckHandler) GetDeckStats(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的卡包ID"))
		return
	}

	stats, err := h.deckService.GetDeckStats(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "获取卡包统计失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(map[string]interface{}{
		"stats": stats,
	}))
}