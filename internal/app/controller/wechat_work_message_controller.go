package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"message-service/internal/app/domain/incoming_http"
	"message-service/internal/app/service"
	"message-service/internal/pkg/logger"
	"net/http"
)

var wwsValidate = validator.New()

func AddWechatWorkRoutes(rg *gin.RouterGroup) {
	wechatWorkMessageGroup := rg.Group("/wechat_work_message")

	wechatWorkMessageGroup.GET("/validation_url", func(context *gin.Context) {
		//	msgSignature string, timestamp string, nonce string, echoStr string
		msgSignature := context.Query("msg_signature")
		timestamp := context.Query("timestamp")
		nonce := context.Query("nonce")
		echoStr := context.Query("echostr")

		logger.Infof("Receving wechat work callback, msg_signature=%s, timestamp=%s, nonce=%s, echostr=%s", msgSignature, timestamp, nonce, echoStr)

		decryptedEchoStr := service.WechatWorkMessageSer.ValidationUrl(msgSignature, timestamp, nonce, echoStr)

		context.String(http.StatusOK, decryptedEchoStr)
	})

	// 异步发送微信消息
	wechatWorkMessageGroup.POST("/", func(context *gin.Context) {
		var createWechatWorkMessageReq incoming_http.CreateWechatWorkMessageReq
		if err := context.ShouldBindJSON(&createWechatWorkMessageReq); err != nil {
			logger.Errorf("Error parsing createWechatWorkMessageReq: %v", err)
			context.JSON(http.StatusOK, Failed())
			return
		}

		if err := wwsValidate.Struct(createWechatWorkMessageReq); err != nil {
			logger.Errorf("Error validating createWechatWorkMessageReq: %v", err)
			context.JSON(http.StatusOK, Failed())
			return
		}

		wechatWorkMessageId, err := service.WechatWorkMessageSer.AddWechatWorkMessage(createWechatWorkMessageReq)
		if err != nil {
			context.JSON(http.StatusOK, Failed())
		} else {
			context.JSON(http.StatusOK, Success(wechatWorkMessageId))
		}
	})
}
