package handlers

import (
	"flashcard/internal/models"
	"flashcard/pkg/database"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testDB *gorm.DB

// setupTestDB 设置测试数据库
func setupTestDB() *gorm.DB {
	if testDB != nil {
		// 清理数据库
		testDB.Exec("DELETE FROM cards")
		testDB.Exec("DELETE FROM tags")
		testDB.Exec("DELETE FROM decks")
		return testDB
	}

	var err error
	testDB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移
	err = testDB.AutoMigrate(&models.Deck{}, &models.Tag{}, &models.Card{}, &models.Review{})
	if err != nil {
		panic("failed to migrate database")
	}

	// 临时替换全局数据库连接
	database.DB = testDB
	
	// 设置一个清理函数，在测试完成后恢复原始数据库连接
	return testDB
}

// setupRouter 设置测试路由
func setupRouter(_ *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// 设置处理器
	deckHandler := NewDeckHandler()
	tagHandler := NewTagHandler()
	cardHandler := NewCardHandler()
	importExportHandler := NewImportExportHandler()

	// 注册路由
	api := r.Group("/api/v1")
	{
		// 卡包路由
		decks := api.Group("/decks")
		{
			decks.GET("", deckHandler.GetDecks)
			decks.POST("", deckHandler.CreateDeck)
			decks.GET("/:id", deckHandler.GetDeck)
			decks.PATCH("/:id", deckHandler.UpdateDeck)
			decks.DELETE("/:id", deckHandler.DeleteDeck)
			decks.GET("/:id/stats", deckHandler.GetDeckStats)
		}

		// 标签路由
		tags := api.Group("/tags")
		{
			tags.GET("/deck/:deckId", tagHandler.GetTags)           // 获取卡包下的所有标签
			tags.POST("/deck/:deckId", tagHandler.CreateTag)         // 创建标签
			tags.GET("/:id", tagHandler.GetTag)         // 获取单个标签
			tags.PATCH("/:id", tagHandler.UpdateTag)    // 更新标签
			tags.DELETE("/:id", tagHandler.DeleteTag)   // 删除标签
			tags.GET("/:id/stats", tagHandler.GetTagStats) // 获取标签统计
			tags.GET("/:id/cards", cardHandler.GetCardsByTag)      // 获取标签下的所有卡片
		}

		// 卡片路由
		cards := api.Group("/cards")
		{
			cards.GET("", cardHandler.SearchCards)
			cards.POST("", cardHandler.CreateCard)
			cards.GET("/:id", cardHandler.GetCard)
			cards.PATCH("/:id", cardHandler.UpdateCard)
			cards.DELETE("/:id", cardHandler.DeleteCard)
		}

		// 导入导出路由
		importExport := api.Group("/import-export")
		{
			importExport.POST("/decks", importExportHandler.ImportDeck)
			importExport.GET("/decks/:deckId", importExportHandler.ExportDeck)
		}
	}

	return r
}