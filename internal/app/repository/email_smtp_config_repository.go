package repository

import (
	"message-service/internal/pkg/db"
	"time"
)

type EmailSmtpConfig struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;column:id"`
	Host      string    `gorm:"type:varchar(255);not null;comment:'host'"`
	Port      int       `gorm:"type:int;not null;comment:'port'"`
	Username  string    `gorm:"type:varchar(255);not null;comment:'用户名'"`
	Password  string    `gorm:"type:varchar(255);not null;comment:'密码'"`
	Status    bool      `gorm:"type:tinyint(1);not null;default:0;comment:'关闭还是开启'"`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP;comment:'更新时间'"`
}

type EmailSmtpConfigRepository struct{}

var EmailSmtpConfigRepo EmailSmtpConfigRepository = EmailSmtpConfigRepository{}

func (*EmailSmtpConfigRepository) FindById(id int64) (*EmailSmtpConfig, error) {
	var emailSmtpConfig EmailSmtpConfig
	if err := db.Db.First(&emailSmtpConfig, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &emailSmtpConfig, nil
}
