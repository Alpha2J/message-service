package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"message-service/internal/app/controller"
	"message-service/internal/pkg/config"
	"message-service/internal/pkg/logger"
	"net/http"
)

func InitHttp() {
	logsDir := logger.GetLogsDir()
	logWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logsDir + "/gin.log",
		MaxSize:    200, // megabytes
		MaxBackups: 10,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})

	if config.IsReleaseMode() {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(gin.LoggerWithWriter(logWriteSyncer))

	// register custom middleware
	r.Use(IncomingRequestLoggingMiddleware())
	r.Use(OwnerValidationMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	rgV1 := r.Group("/v1")
	controller.AddEmailRoutes(rgV1)
	controller.AddWechatWorkRoutes(rgV1)

	port := fmt.Sprintf("%d", config.GetAppConfig().Port)
	logger.Info("HTTP server started with port :" + port)
	err := r.Run(":" + port)
	if err != nil {
		return
	}
}
