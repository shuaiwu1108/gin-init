package main

import (
	"fmt"
	"gin-init/config"
	"gin-init/db"
	_ "gin-init/docs"
	"gin-init/logger"
	"gin-init/router"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:generate swag init

// @title           Gin init API
// @version         1.0
// @description     Gin-Init-API

func main() {
	// 创建Gin引擎
	r := gin.Default()

	// 配置 CORS 中间件，允许跨域请求
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                      // 开发环境可临时用 * 允许所有域名，生产环境需指定具体域名
		AllowMethods:     []string{"GET", "POST", "OPTIONS"}, // 允许的请求方法（Swagger 主要用 GET）
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", base)

	// 配置Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 模块
	{
		develop := r.Group("/api")

		// 搜索模块
		{
			search := develop.Group("/v1")
			search.GET("/test", router.Test)
		}
	}
	log.Println(fmt.Sprintf("接口地址：http://127.0.0.1:%d", config.Cfg.App.Port))
	log.Println(fmt.Sprintf("接口文档地址：http://127.0.0.1:%d/swagger/index.html", config.Cfg.App.Port))
	// 启动服务
	err := r.Run(fmt.Sprintf(":%d", config.Cfg.App.Port))
	if err != nil {
		return
	}
}

func init() {
	// 步骤1：加载配置文件
	if err := config.Init("app.yaml"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 步骤2：初始化日志切割（在 Gin 引擎创建前）
	logger.Init()

	// 步骤3：初始化 GORM（依赖配置，在 Gin 引擎创建前执行）
	if err := db.Init(); err != nil {
		log.Fatalf("GORM 初始化失败: %v", err)
	}

	// 根据配置设置Gin模式
	if config.Cfg.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

func base(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Gin init is started!",
	})
}
