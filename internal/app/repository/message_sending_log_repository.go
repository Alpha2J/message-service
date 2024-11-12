package repository

import (
	"message-service/internal/pkg/db"
	"time"
)

type MessageSendingLog struct {
	ID           int64     `gorm:"primaryKey;autoIncrement;column:id"`
	MessageID    int64     `gorm:"not null;column:message_id;comment:'邮件id'"`
	IsSuccess    bool      `gorm:"type:tinyint(1);not null;comment:'是否成功'"`
	FailedReason string    `gorm:"type:text;not null;comment:'失败原因'"`
	Type         string    `gorm:"type:varchar(255);not null;comment:'消息类型，邮件还是微信消息'"`
	CreatedAt    time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
	UpdatedAt    time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP;comment:'更新时间'"`
}

const (
	MESSAGE_SENDING_LOG_TYPE_EMAIL               string = "EMAIL"
	MESSAGE_SENDING_LOG_TYPE_WECHAT_WORK_MESSAGE string = "WECHAT_WORK_MESSAGE"
)

type MessageSendingLogRepository struct{}

var MessageSendingLogRepo MessageSendingLogRepository = MessageSendingLogRepository{}

func (*MessageSendingLogRepository) Save(messageSendingLog *MessageSendingLog) error {
	result := db.Db.Create(messageSendingLog)
	return result.Error
}
