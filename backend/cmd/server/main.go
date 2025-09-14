package main

import (
	"log"
	"net/http"

	"flashcard/internal/config"
	"flashcard/internal/handlers"
	"flashcard/internal/middleware"
	"flashcard/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 设置Gin模式
	gin.SetMode(cfg.GinMode)

	// 初始化数据库
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.CloseDB()

	// 创建处理器
	deckHandler := handlers.NewDeckHandler()
	tagHandler := handlers.NewTagHandler()
	cardHandler := handlers.NewCardHandler()
	importExportHandler := handlers.NewImportExportHandler()
	studyHandler := handlers.NewStudyHandler()
	systemHandler := handlers.NewSystemHandler()

	// 创建Gin引擎
	r := gin.New()

	// 使用中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Flashcard API 服务运行正常",
		})
	})

	// API路由组
	api := r.Group("/api/v1")
	{
		// 临时测试接口
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// 卡包相关路由
		decks := api.Group("/decks")
		{
			decks.GET("", deckHandler.GetDecks)                 // 获取所有卡包
			decks.POST("", deckHandler.CreateDeck)              // 创建卡包
			decks.GET("/:id", deckHandler.GetDeck)              // 获取单个卡包
			decks.PATCH("/:id", deckHandler.UpdateDeck)         // 更新卡包
			decks.DELETE("/:id", deckHandler.DeleteDeck)        // 删除卡包
			decks.GET("/:id/stats", deckHandler.GetDeckStats)   // 获取卡包统计
			decks.GET("/:id/cards", cardHandler.GetCardsByDeck) // 获取卡包下的所有卡片
		}

		// 标签相关路由
		apiTags := api.Group("/tags")
		{
			apiTags.GET("", tagHandler.GetAllTags)               // 获取所有标签
			apiTags.POST("", tagHandler.CreateTag)               // 创建标签
			apiTags.GET("/deck/:deckId", tagHandler.GetTags)     // 获取卡包下的所有标签
			apiTags.GET("/:id", tagHandler.GetTag)               // 获取单个标签
			apiTags.PATCH("/:id", tagHandler.UpdateTag)          // 更新标签
			apiTags.DELETE("/:id", tagHandler.DeleteTag)         // 删除标签
			apiTags.GET("/:id/stats", tagHandler.GetTagStats)    // 获取标签统计
			apiTags.GET("/:id/cards", cardHandler.GetCardsByTag) // 获取标签下的所有卡片
		}

		// 卡片相关路由
		apiCards := api.Group("/cards")
		{
			apiCards.GET("", cardHandler.SearchCards)       // 搜索卡片
			apiCards.POST("", cardHandler.CreateCard)       // 创建卡片
			apiCards.GET("/:id", cardHandler.GetCard)       // 获取单个卡片
			apiCards.PATCH("/:id", cardHandler.UpdateCard)  // 更新卡片
			apiCards.DELETE("/:id", cardHandler.DeleteCard) // 删除卡片
		}

		// 导入导出相关路由
		apiImportExport := api.Group("/import-export")
		{
			apiImportExport.POST("/decks", importExportHandler.ImportDeck)        // 导入卡包
			apiImportExport.GET("/decks/:deckId", importExportHandler.ExportDeck) // 导出卡包
		}

		// 学习相关路由
		apiStudy := api.Group("/study")
		{
			apiStudy.POST("/deck/:deckId", studyHandler.StartDeckStudy) // 开始学习卡包
			apiStudy.POST("/tag/:tagId", studyHandler.StartTagStudy)    // 开始学习标签
			apiStudy.POST("/random", studyHandler.StartRandomStudy)     // 开始随机学习
			apiStudy.GET("/due", studyHandler.GetDueCards)              // 获取到期卡片
			apiStudy.POST("/review/:cardId", studyHandler.SubmitReview) // 提交复习结果
		}

		// 系统管理相关路由
		apiSystem := api.Group("/system")
		{
			apiSystem.GET("/backup", systemHandler.BackupData)     // 备份数据
			apiSystem.POST("/restore", systemHandler.RestoreData)  // 恢复数据
			apiSystem.DELETE("/clear", systemHandler.ClearAllData) // 清空数据
			apiSystem.GET("/stats", systemHandler.GetSystemStats)  // 获取系统统计
		}
	}

	// 启动服务器
	log.Printf("服务器启动在端口 %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
