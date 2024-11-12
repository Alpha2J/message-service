package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"message-service/internal/app/domain/incoming_http"
	"message-service/internal/app/service"
	"message-service/internal/pkg/logger"
	"net/http"
)

var eValidate = validator.New()

func AddEmailRoutes(rg *gin.RouterGroup) {
	emailGroup := rg.Group("/email")

	// 异步发送邮件
	emailGroup.POST("/", func(context *gin.Context) {
		var createEmailReq incoming_http.CreateEmailReq
		if err := context.ShouldBindJSON(&createEmailReq); err != nil {
			logger.Errorf("Error parsing createEmailReq: %v", err)
			context.JSON(http.StatusOK, Failed())
			return
		}

		if err := eValidate.Struct(createEmailReq); err != nil {
			logger.Errorf("Error validating createEmailReq: %v", err)
			context.JSON(http.StatusOK, Failed())
			return
		}

		emailId, err := service.EmailSer.AddEmail(createEmailReq)
		if err != nil {
			context.JSON(http.StatusOK, Failed())
		} else {
			context.JSON(http.StatusOK, Success(emailId))
		}
	})

	// 获取单个邮件的详细信息
	emailGroup.GET("/:emailId", func(context *gin.Context) {
		emailId := context.Param("emailId")

		err := eValidate.Var(emailId, "required")
		if err != nil {
			logger.Errorf("Error validating emailId: %v", err)
			context.JSON(http.StatusOK, Failed())
			return
		}

		emailVO, err := service.EmailSer.GetEmail(emailId)
		if err != nil {
			context.JSON(http.StatusOK, Failed())
		} else {
			context.JSON(http.StatusOK, Success(emailVO))
		}
	})
}
