package app

import (
	"github.com/gin-gonic/gin"
	"message-service/internal/pkg/config"
	"message-service/internal/pkg/logger"
	"net/http"
)

var unrestrictedPaths = []string{"/ping", "/v1/wechat_work_message/validation_url"}

func isUnrestrictedPath(path string) bool {
	for _, unrestrictedPath := range unrestrictedPaths {
		if path == unrestrictedPath {
			return true
		}
	}
	return false
}

func OwnerValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 白名单内请求直接走业务逻辑，不经过middleware处理
		if isUnrestrictedPath(c.FullPath()) {
			c.Next()
			return
		}

		// 获取请求头中的特定字段（例如 "owner-token"）
		headerValue := c.GetHeader("owner-token")

		ownerToken := config.GetBusinessConfig().OwnerToken
		// 如果请求头中没有该字段或字段值为空，返回错误
		if headerValue != ownerToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized request!!!",
			})
			// 终止请求
			c.Abort()
			return
		}

		c.Next()
	}
}

func IncomingRequestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		url := c.Request.URL.String()
		logger.Infof("Incomming Request, method: %s, url: %s", method, url)

		c.Next()
	}
}
