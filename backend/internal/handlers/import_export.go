package handlers

import (
	"flashcard/internal/models"
	"flashcard/internal/services"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ImportExportHandler 导入导出处理器
type ImportExportHandler struct {
	importExportService *services.ImportExportService
}

// NewImportExportHandler 创建导入导出处理器实例
func NewImportExportHandler() *ImportExportHandler {
	return &ImportExportHandler{
		importExportService: services.NewImportExportService(),
	}
}

// ExportDeck 导出卡包
func (h *ImportExportHandler) ExportDeck(c *gin.Context) {
	deckIDStr := c.Param("deckId")
	deckID, err := strconv.ParseUint(deckIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "无效的卡包ID"))
		return
	}

	// 获取导出格式
	format := c.DefaultQuery("format", "json")

	// 导出卡包
	filename, err := h.importExportService.ExportDeck(uint(deckID), format)
	if err != nil {
		// 检查是否是记录不存在的错误
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse(models.CodeNotFound, "导出卡包失败", "卡包不存在"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "导出卡包失败", err.Error()))
		return
	}

	// 检查查询参数决定是下载文件还是返回JSON
	download := c.DefaultQuery("download", "false")
	if download == "true" {
		// 直接下载文件
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename="+filepath.Base(filename))
		c.Header("Content-Type", "application/octet-stream")
		c.File(filename)

		// 清理临时文件
		go func() {
			os.Remove(filename)
		}()
		return
	}

	// 返回JSON响应（包含文件名）
	response := models.SuccessResponse(gin.H{
		"filename":     filename,
		"download_url": "/api/v1/import-export/decks/" + deckIDStr + "?download=true&format=" + format,
	})
	response.Message = "卡包导出成功"
	c.JSON(http.StatusOK, response)
}

// ImportDeck 导入卡包
func (h *ImportExportHandler) ImportDeck(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "请选择要导入的文件"))
		return
	}

	// 检查文件大小（限制为10MB）
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "文件大小不能超过10MB"))
		return
	}

	// 创建临时文件
	tempDir := "./temp"
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		os.Mkdir(tempDir, 0755)
	}

	tempFile, err := os.CreateTemp(tempDir, "upload-*"+filepath.Ext(file.Filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "创建临时文件失败", err.Error()))
		return
	}
	defer os.Remove(tempFile.Name())

	// 保存上传的文件到临时文件
	if err := c.SaveUploadedFile(file, tempFile.Name()); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "保存文件失败", err.Error()))
		return
	}

	// 导入卡包
	deck, err := h.importExportService.ImportDeck(tempFile.Name())
	if err != nil {
		// 检查是否是不支持的格式错误
		if strings.Contains(err.Error(), "不支持的导入格式") {
			c.JSON(http.StatusBadRequest, models.ErrorResponse(models.CodeInvalidParam, "导入失败", "不支持的文件格式"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.CodeInternal, "导入卡包失败", err.Error()))
		return
	}

	// 获取导入的卡片数量
	var cardCount int64
	h.importExportService.GetDB().Model(&models.Card{}).Where("deck_id = ?", deck.ID).Count(&cardCount)

	response := models.SuccessResponse(gin.H{
		"deck_id":     deck.ID,
		"deck_name":   deck.Name,
		"card_count":  cardCount,
		"import_time": deck.CreatedAt,
	})
	response.Message = fmt.Sprintf("成功导入卡包\"%s\"，包含 %d 张卡片", deck.Name, cardCount)
	c.JSON(http.StatusCreated, response)
}
