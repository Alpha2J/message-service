package service

import (
	"gopkg.in/gomail.v2"
	"message-service/internal/app/domain/incoming_http"
	"message-service/internal/app/manager"
	"message-service/internal/app/repository"
	"message-service/internal/pkg/helper"
	"message-service/internal/pkg/logger"
	"sync"
	"time"
)

type EmailService struct{}

var EmailSer EmailService = EmailService{}

func (*EmailService) AddEmail(createEmailReq incoming_http.CreateEmailReq) (string, error) {
	logger.Infof("Adding email, req=%v", createEmailReq)

	emailId := helper.NewUUIDStr()
	jsonStr, _ := helper.StringArrToJSONStr(createEmailReq.To)
	emailEntity := &repository.Email{
		UUID:    emailId,
		From:    createEmailReq.From,
		To:      jsonStr,
		Subject: createEmailReq.Subject,
		Body:    createEmailReq.Body,
		Status:  repository.EMAIL_STATUS_INIT,
	}

	err := repository.EmailRepo.Save(emailEntity)
	if err != nil {
		return "", err
	}

	return emailId, nil
}

func (*EmailService) GetEmail(uuid string) (*incoming_http.EmailVO, error) {
	emailEntity, err := repository.EmailRepo.FindByUUID(uuid)
	if err != nil {
		return nil, err
	}

	stringArr, _ := helper.JSONStrToStringArr(emailEntity.To)
	return &incoming_http.EmailVO{
		Id:      emailEntity.UUID,
		From:    emailEntity.From,
		To:      stringArr,
		Subject: emailEntity.Subject,
		Body:    emailEntity.Body,
	}, nil
}

var isEmailSendingTaskLocked bool
var emailSendingTaskMu sync.Mutex

func (*EmailService) EmailSendingTask() error {
	if isEmailSendingTaskLocked {
		logger.Info("Still processing the email sending task.")
		return nil
	}
	emailSendingTaskMu.Lock()
	isEmailSendingTaskLocked = true
	defer func() {
		isEmailSendingTaskLocked = false
		emailSendingTaskMu.Unlock()
	}()

	sendingEmails, err := repository.EmailRepo.FindPendingEmails()
	if err != nil {
		logger.Errorf("Get sending emails err: %v", err)
		return err
	}

	if len(sendingEmails) == 0 {
		return nil
	}

	maxConcurrent := 2
	// should we change repository.Email to *repository.Email
	ch := make(chan repository.Email, maxConcurrent)

	for _, sendingEmail := range sendingEmails {
		err := repository.EmailRepo.UpdateStatusById(sendingEmail.ID, repository.EMAIL_STATUS_SENDING)
		if err != nil {
			continue
		}

		ch <- sendingEmail
		go sendEmail(ch)
	}

	return nil
}

func sendEmail(ch chan repository.Email) error {
	email := <-ch

	smtpConfig, err := manager.EmailSmtpConfigMana.FindByEmailSenderAddress(email.From)
	if err != nil {
		return err
	}

	smtpHost := smtpConfig.Host
	smtpPort := smtpConfig.Port
	smtpUsername := smtpConfig.Username
	smtpPassword := smtpConfig.Password

	from := email.From
	to, _ := helper.JSONStrToStringArr(email.To)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/plain", email.Body)

	// 设置HTML格式的邮件正文
	//m.AddAlternative("text/html", "<p>This is a test email sent using <b>gomail</b> in Go.</p>")
	// 添加附件（可选）
	//m.Attach("/path/to/file.txt")

	// 设置SMTP服务器信息
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		// set resend_count and resend_at and email status to failed
		// add email_sending_log
		resendAt := time.Now().Add(5 * time.Second)
		repository.EmailRepo.UpdateStatusAndResendCountAndResendAtById(email.ID, repository.EMAIL_STATUS_FAILED, email.ResendCount+1, &resendAt)
		repository.MessageSendingLogRepo.Save(&repository.MessageSendingLog{
			MessageID: email.ID,
			IsSuccess: false,
			Type:      repository.MESSAGE_SENDING_LOG_TYPE_EMAIL,
		})
		return err
	} else {
		repository.EmailRepo.UpdateStatusById(email.ID, repository.EMAIL_STATUS_SUCCESS)
		repository.MessageSendingLogRepo.Save(&repository.MessageSendingLog{
			MessageID: email.ID,
			IsSuccess: true,
			Type:      repository.MESSAGE_SENDING_LOG_TYPE_EMAIL,
		})
		return nil
	}
}
