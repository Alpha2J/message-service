package repository

import (
	"message-service/internal/pkg/db"
	"time"
)

type Email struct {
	ID          int64      `gorm:"primaryKey;autoIncrement;column:id"`
	UUID        string     `gorm:"type:varchar(255);not null;comment:'唯一id'"`
	From        string     `gorm:"type:varchar(255);not null;comment:'发送者'"`
	To          string     `gorm:"type:varchar(1024);not null;comment:'接受者，逗号隔开多个接受者'"`
	Subject     string     `gorm:"type:varchar(1024);not null;comment:'主题'"`
	Body        string     `gorm:"type:text;not null;comment:'邮件内容'"`
	Status      string     `gorm:"type:varchar(255);not null;comment:'状态'"`
	ResendCount int        `gorm:"type:int;not null;default:0;comment:'重发次数'"`
	ResendAt    *time.Time `gorm:"type:datetime;comment:'下次重发时间'"`
	CreatedAt   time.Time  `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
	UpdatedAt   time.Time  `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP;comment:'更新时间'"`
}

const (
	EMAIL_STATUS_INIT    string = "INIT"
	EMAIL_STATUS_SENDING string = "SENDING"
	EMAIL_STATUS_SUCCESS string = "SUCCESS"
	EMAIL_STATUS_FAILED  string = "FAILED"
)

type EmailRepository struct{}

var EmailRepo EmailRepository = EmailRepository{}

func (*EmailRepository) Save(email *Email) error {
	result := db.Db.Create(email)
	return result.Error
}

func (*EmailRepository) FindByUUID(uuid string) (*Email, error) {
	var email Email
	if err := db.Db.First(&email, "uuid = ?", uuid).Error; err != nil {
		return nil, err
	}

	return &email, nil
}

func (*EmailRepository) FindAllByStatus(status string) ([]Email, error) {
	var emails []Email
	if err := db.Db.Where("status = ?", status).Find(&emails).Error; err != nil {
		return nil, err
	}

	return emails, nil
}

func (*EmailRepository) UpdateStatusById(id int64, status string) error {
	if err := db.Db.Model(&Email{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

func (*EmailRepository) UpdateStatusAndResendCountAndResendAtById(id int64, status string, resendCount int, resendAt *time.Time) error {
	if err := db.Db.Model(&Email{}).
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

func (*EmailRepository) FindPendingEmails() ([]Email, error) {
	var emails []Email
	if err := db.Db.Where("status = ?", EMAIL_STATUS_INIT).Or("status = ? AND resend_count < ? AND resend_at <= ?", EMAIL_STATUS_FAILED, 5, time.Now()).Find(&emails).Error; err != nil {
		return nil, err
	}

	return emails, nil
}
