package manager

import (
	"message-service/internal/app/repository"
	"message-service/internal/pkg/helper"
)

type EmailSmtpConfigManager struct{}

var EmailSmtpConfigMana EmailSmtpConfigManager = EmailSmtpConfigManager{}

var emailSmtpConfigCache *helper.Cache = helper.NewCache()

func (*EmailSmtpConfigManager) FindByEmailSenderAddress(address string) (*repository.EmailSmtpConfig, error) {
	value, found := emailSmtpConfigCache.Get(address)
	if found {
		return value.(*repository.EmailSmtpConfig), nil
	}

	emailSender, err := repository.EmailSenderRepo.FindByAddress(address)
	if err != nil {
		return nil, err
	}

	smtpConfig, err := repository.EmailSmtpConfigRepo.FindById(emailSender.EmailSmtpConfigID)
	if err != nil {
		return nil, err
	}

	emailSmtpConfigCache.Set(address, smtpConfig)

	return smtpConfig, nil
}
