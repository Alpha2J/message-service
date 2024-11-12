package app

import (
	"github.com/robfig/cron/v3"
	"message-service/internal/app/service"
	"message-service/internal/pkg/logger"
)

func InitScheduler() {
	c := cron.New(cron.WithSeconds())

	logger.Info("Scheduler cron started.")

	c.AddFunc("*/2 * * * * *", func() {
		logger.Info("Executing email sending task.")
		_ = service.EmailSer.EmailSendingTask()
	})

	c.AddFunc("*/2 * * * * *", func() {
		logger.Info("Executing wechat work message sending task.")
		_ = service.WechatWorkMessageSer.WechatWorkMessageSendingTask()
	})

	// 启动定时任务
	c.Start()
}
