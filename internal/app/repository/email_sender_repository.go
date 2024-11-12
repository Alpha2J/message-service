package repository

import (
	"message-service/internal/pkg/db"
	"time"
)

type EmailSender struct {
	ID                int64     `gorm:"primaryKey;autoIncrement;column:id"`
	EmailSmtpConfigID int64     `gorm:"not null;column:email_smtp_config_id"`
	Address           string    `gorm:"type:varchar(255);not null;uniqueIndex;comment:'email address'"`
	CreatedAt         time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
	UpdatedAt         time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP;comment:'更新时间'"`
}

type EmailSenderRepository struct{}

var EmailSenderRepo EmailSenderRepository = EmailSenderRepository{}

func (*EmailSenderRepository) FindByAddress(address string) (*EmailSender, error) {
	var emailSender EmailSender
	if err := db.Db.First(&emailSender, "address = ?", address).Error; err != nil {
		return nil, err
	}

	return &emailSender, nil
}
