package repository

import (
	"message-service/internal/pkg/db"
	"time"
)

type WechatWorkAppConfig struct {
	ID                       int64     `gorm:"primaryKey;autoIncrement;column:id"`
	CompanyID                string    `gorm:"type:varchar(255);not null;comment:'应用所属的企业id'"`
	CompanySecret            string    `gorm:"type:varchar(1024);not null;comment:'应用所属的企业id的secret'"`
	AgentID                  string    `gorm:"type:varchar(255);not null;comment:'企微那边的app id'"`
	MsgReceivingServerToken  string    `gorm:"type:varchar(255);not null;comment:'接收消息服务器的token'"`
	MsgReceivingServerAeskey string    `gorm:"type:varchar(255);not null;comment:'接收消息服务器的aeskey'"`
	Description              *string   `gorm:"type:varchar(1024);comment:'描述'"`
	CreatedAt                time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
	UpdatedAt                time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP;comment:'更新时间'"`
}

type WechatWorkAppConfigRepository struct{}

var WechatWorkAppConfigRepo WechatWorkAppConfigRepository = WechatWorkAppConfigRepository{}

func (*WechatWorkAppConfigRepository) FindById(id int64) (*WechatWorkAppConfig, error) {
	var wechatWorkAppConfig WechatWorkAppConfig
	if err := db.Db.First(&wechatWorkAppConfig, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &wechatWorkAppConfig, nil
}
