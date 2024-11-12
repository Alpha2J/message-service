package repository

import (
	"message-service/internal/pkg/db"
	"time"
)

type WechatWorkMessage struct {
	ID              int64      `gorm:"primaryKey;autoIncrement;column:id"`
	UUID            string     `gorm:"type:varchar(255);not null;comment:'唯一id'"`
	FromAppConfigID int64      `gorm:"not null;comment:'发送者, 使用哪个应用发送, 对应wechat_work_app_configs的id字段'"`
	To              string     `gorm:"type:text;not null;comment:'接受者，逗号隔开多个接受者'"`
	Content         string     `gorm:"type:text;not null;comment:'消息内容'"`
	Status          string     `gorm:"type:varchar(255);not null;comment:'状态'"`
	ResendCount     int        `gorm:"type:int;not null;default:0;comment:'重发次数'"`
	ResendAt        *time.Time `gorm:"type:datetime;comment:'下次重发时间'"`
	CreatedAt       time.Time  `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
	UpdatedAt       time.Time  `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP;comment:'更新时间'"`
}

const (
	WECHAT_WORK_MESSAGE_STATUS_INIT    string = "INIT"
	WECHAT_WORK_MESSAGE_STATUS_SENDING string = "SENDING"
	WECHAT_WORK_MESSAGE_STATUS_SUCCESS string = "SUCCESS"
	WECHAT_WORK_MESSAGE_STATUS_FAILED  string = "FAILED"
)

type WechatWorkMessageRepository struct{}

var WechatWorkMessageRepo WechatWorkMessageRepository = WechatWorkMessageRepository{}

func (*WechatWorkMessageRepository) Save(wechatWorkMessage *WechatWorkMessage) error {
	result := db.Db.Create(wechatWorkMessage)
	return result.Error
}

func (*WechatWorkMessageRepository) FindPendingWechatWorkMessages() ([]WechatWorkMessage, error) {
	var wechatWorkMessages []WechatWorkMessage
	if err := db.Db.Where("status = ?", WECHAT_WORK_MESSAGE_STATUS_INIT).Or("status = ? AND resend_count < ? AND resend_at <= ?", WECHAT_WORK_MESSAGE_STATUS_FAILED, 5, time.Now()).Find(&wechatWorkMessages).Error; err != nil {
		return nil, err
	}

	return wechatWorkMessages, nil
}

func (*WechatWorkMessageRepository) UpdateStatusById(id int64, status string) error {
	if err := db.Db.Model(&WechatWorkMessage{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

func (*WechatWorkMessageRepository) UpdateStatusAndResendCountAndResendAtById(id int64, status string, resendCount int, resendAt *time.Time) error {
	if err := db.Db.Model(&WechatWorkMessage{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"resend_count": resendCount,
			"resend_at":    resendAt,
		}).Error; err != nil {
		return err
	}

	return nil
}
