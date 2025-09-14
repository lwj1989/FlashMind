package handlers

import (
	"flashcard/internal/models"
	"flashcard/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// StudyHandler 学习处理器
type StudyHandler struct {
	studyService *services.StudyService
}

// NewStudyHandler 创建学习处理器实例
func NewStudyHandler() *StudyHandler {
	return &StudyHandler{
		studyService: services.NewStudyService(),
	}
}

// StartDeckStudy 开始学习卡包
func (h *StudyHandler) StartDeckStudy(c *gin.Context) {
	deckIDStr := c.Param("deckId")
	deckID, err := strconv.ParseUint(deckIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的卡包ID"))
		return
	}

	// 获取学习队列限制
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	session, err := h.studyService.StartDeckStudy(uint(deckID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "开始学习失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(session))
}

// StartTagStudy 开始学习标签
func (h *StudyHandler) StartTagStudy(c *gin.Context) {
	tagIDStr := c.Param("tagId")
	tagID, err := strconv.ParseUint(tagIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的标签ID"))
		return
	}

	// 获取学习队列限制
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	session, err := h.studyService.StartTagStudy(uint(tagID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "开始学习失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(session))
}

// StartRandomStudy 开始随机学习
func (h *StudyHandler) StartRandomStudy(c *gin.Context) {
	// 获取学习队列限制
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	session, err := h.studyService.StartRandomStudy(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "开始随机学习失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(session))
}

// SubmitReview 提交复习结果
func (h *StudyHandler) SubmitReview(c *gin.Context) {
	cardIDStr := c.Param("cardId")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的卡片ID"))
		return
	}

	var req models.ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "请求参数格式错误: "+err.Error()))
		return
	}

	// 验证 ReviewResult 范围
	if req.Result < 0 || req.Result > 2 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "复习结果必须在0-2之间"))
		return
	}

	response, err := h.studyService.SubmitReview(uint(cardID), req.Result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "提交复习结果失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(response))
}

// GetDueCards 获取到期复习的卡片
func (h *StudyHandler) GetDueCards(c *gin.Context) {
	// 获取限制
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	session, err := h.studyService.GetDueCards(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "获取到期卡片失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(session))
}
