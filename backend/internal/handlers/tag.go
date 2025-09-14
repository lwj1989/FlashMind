package handlers

import (
	"flashcard/internal/models"
	"flashcard/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TagHandler 标签处理器
type TagHandler struct {
	tagService *services.TagService
}

// NewTagHandler 创建标签处理器实例
func NewTagHandler() *TagHandler {
	return &TagHandler{
		tagService: services.NewTagService(),
	}
}

// CreateTag 创建标签
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required,min=1,max=50"`
		DeckID *uint  `json:"deck_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "标签名称不能为空且长度应在1-50之间"))
		return
	}

	tag, err := h.tagService.CreateTag(req.DeckID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "创建标签失败", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse(tag))
}

// GetTags 获取卡包下的所有标签
func (h *TagHandler) GetTags(c *gin.Context) {
	deckIDStr := c.Param("deckId")
	if deckIDStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "卡包ID不能为空"))
		return
	}
	deckID, err := strconv.ParseUint(deckIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的卡包ID"))
		return
	}

	includeStats := c.DefaultQuery("include_stats", "false")

	if includeStats == "true" {
		tags, err := h.tagService.GetTagsByDeckIDWithStats(uint(deckID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "获取标签列表失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse(tags))
	} else {
		tags, err := h.tagService.GetTagsByDeckID(uint(deckID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "获取标签列表失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse(tags))
	}
}

// GetTag 获取单个标签
func (h *TagHandler) GetTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的标签ID"))
		return
	}

	tag, err := h.tagService.GetTagByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(models.CodeNotFound, "标签不存在"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(tag))
}

// UpdateTag 更新标签
func (h *TagHandler) UpdateTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的标签ID"))
		return
	}

	var req struct {
		Name   string `json:"name" binding:"required,min=1,max=50"`
		DeckID *uint  `json:"deck_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "标签名称不能为空且长度应在1-50之间"))
		return
	}

	tag, err := h.tagService.UpdateTagWithDeck(uint(id), req.Name, req.DeckID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "更新标签失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(tag))
}

// DeleteTag 删除标签
func (h *TagHandler) DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的标签ID"))
		return
	}

	// 获取查询参数，决定是否同时删除卡片
	deleteCardsStr := c.DefaultQuery("delete_cards", "false")
	deleteCards := deleteCardsStr == "true"

	if err := h.tagService.DeleteTag(uint(id), deleteCards); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "删除标签失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(nil))
}

// GetAllTags 获取所有标签
func (h *TagHandler) GetAllTags(c *gin.Context) {
	includeStats := c.DefaultQuery("include_stats", "false")

	if includeStats == "true" {
		tags, err := h.tagService.GetAllTagsWithStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "获取标签列表失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse(tags))
	} else {
		tags, err := h.tagService.GetAllTags()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "获取标签列表失败", err.Error()))
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse(tags))
	}
}

// GetTagStats 获取标签统计信息
func (h *TagHandler) GetTagStats(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的标签ID"))
		return
	}

	stats, err := h.tagService.GetTagStats(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "获取标签统计失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(stats))
}