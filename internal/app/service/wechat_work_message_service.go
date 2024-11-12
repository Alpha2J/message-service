package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"message-service/internal/app/domain/incoming_http"
	"message-service/internal/app/manager"
	"message-service/internal/app/repository"
	"message-service/internal/app/service/wechat_lib"
	"message-service/internal/pkg/helper"
	"message-service/internal/pkg/logger"
	"net/http"
	"sync"
	"time"
)

type WechatWorkMessageService struct{}

var WechatWorkMessageSer WechatWorkMessageService = WechatWorkMessageService{}

func (*WechatWorkMessageService) ValidationUrl(msgSignature string, timestamp string, nonce string, echoStr string) string {

	wechatWorkAppConfig, err := manager.WechatWorkAppConfigMana.FindById(1)
	if err != nil {
		logger.Errorf("Get wechat work app config err: %v", err)
		return ""
	}

	wxcpt := wechat_lib.NewWXBizMsgCrypt(wechatWorkAppConfig.MsgReceivingServerToken, wechatWorkAppConfig.MsgReceivingServerAeskey, wechatWorkAppConfig.CompanyID, wechat_lib.JsonType)
	decryptedEchoStr, cryptError := wxcpt.VerifyURL(msgSignature, timestamp, nonce, echoStr)
	if nil != cryptError {
		logger.Errorf("Wechat work callback url verifiy failed: %v", err)
		return ""
	}

	logger.Infof("Decrypted wechat work EchoStr: %v", string(decryptedEchoStr))

	return string(decryptedEchoStr)
}

func (*WechatWorkMessageService) AddWechatWorkMessage(createWechatWorkMessageReq incoming_http.CreateWechatWorkMessageReq) (string, error) {
	logger.Infof("Adding wechat work message, req: %v", createWechatWorkMessageReq)

	wechatWorkMessageId := helper.NewUUIDStr()
	jsonStr, _ := helper.StringArrToJSONStr(createWechatWorkMessageReq.To)
	emailEntity := &repository.WechatWorkMessage{
		UUID:            wechatWorkMessageId,
		FromAppConfigID: createWechatWorkMessageReq.FromAppID,
		To:              jsonStr,
		Content:         createWechatWorkMessageReq.Content,
		Status:          repository.WECHAT_WORK_MESSAGE_STATUS_INIT,
	}

	err := repository.WechatWorkMessageRepo.Save(emailEntity)
	if err != nil {
		return "", err
	}

	return wechatWorkMessageId, nil
}

var isWechatWorkMessageSendingTaskLocked bool
var wechatWorkMessageSendingTaskMu sync.Mutex

func (*WechatWorkMessageService) WechatWorkMessageSendingTask() error {
	if isWechatWorkMessageSendingTaskLocked {
		logger.Info("Still processing the wechatWorkMessageSendingTask, continue...")
		return nil
	}
	wechatWorkMessageSendingTaskMu.Lock()
	isWechatWorkMessageSendingTaskLocked = true
	defer func() {
		isWechatWorkMessageSendingTaskLocked = false
		wechatWorkMessageSendingTaskMu.Unlock()
	}()

	sendingWechatWorkMessages, err := repository.WechatWorkMessageRepo.FindPendingWechatWorkMessages()
	if err != nil {
		logger.Errorf("Get sending wechat work messages err: %v", err)
		return err
	}

	if len(sendingWechatWorkMessages) == 0 {
		return nil
	}

	maxConcurrent := 2
	ch := make(chan repository.WechatWorkMessage, maxConcurrent)

	for _, sendingWechatWorkMessage := range sendingWechatWorkMessages {
		err := repository.WechatWorkMessageRepo.UpdateStatusById(sendingWechatWorkMessage.ID, repository.WECHAT_WORK_MESSAGE_STATUS_SENDING)
		if err != nil {
			continue
		}

		ch <- sendingWechatWorkMessage
		go sendWechatWorkMessage(ch)
	}

	return nil
}

func sendWechatWorkMessage(ch chan repository.WechatWorkMessage) error {
	wechatWorkMessage := <-ch

	appConfig, err := manager.WechatWorkAppConfigMana.FindById(wechatWorkMessage.FromAppConfigID)
	if err != nil {
		return err
	}

	corpId := appConfig.CompanyID
	corpSecret := appConfig.CompanySecret

	// get access token
	getAccessTokenUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", corpId, corpSecret)
	getAccessTokenResp, err := http.Get(getAccessTokenUrl)
	if err != nil {
		logger.Errorf("Get access token err: %v", err)
		return err
	}
	defer getAccessTokenResp.Body.Close()

	getAccessTokenRespBody, err := io.ReadAll(getAccessTokenResp.Body)
	if err != nil {
		logger.Errorf("Read response body err: %v", err)
		return err
	}

	var getAccessTokenResponse GetAccessTokenResponse
	err = json.Unmarshal(getAccessTokenRespBody, &getAccessTokenResponse)
	if err != nil {
		logger.Errorf("Unmarshal response body err: %v", err)
		return err
	}

	logger.Infof("Get access token response: %v", getAccessTokenResponse.AccessToken)

	sendMessageUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?debug=1&access_token=%s", getAccessTokenResponse.AccessToken)
	// todo change ToUser
	sendMessageReq := SendWechatWorkTextMessageReq{
		ToUser:  "WangJiaJue",
		MsgType: "text",
		AgentId: appConfig.AgentID,
		Text: SendWechatWorkTextMessage{
			Content: wechatWorkMessage.Content,
		},
	}
	sendMessageReqJson, err := json.Marshal(&sendMessageReq)
	if err != nil {
		logger.Errorf("Marshal json err: %v", err)
		return err
	}

	sendMessageResp, err := http.Post(sendMessageUrl, "application/json", bytes.NewBuffer(sendMessageReqJson))
	if err != nil {
		logger.Errorf("Get access token err: %v", err)
		return err
	}
	defer sendMessageResp.Body.Close()

	sendMessageResponseBody, err := io.ReadAll(sendMessageResp.Body)
	if err != nil {
		logger.Errorf("Read response body err: %v", err)
		return err
	}

	var sendMessageResponse SendWechatWorkTextMessageResp
	err = json.Unmarshal(sendMessageResponseBody, &sendMessageResponse)
	if err != nil {
		logger.Errorf("Unmarshal response body err: %v", err)
		return err
	}

	if sendMessageResponse.ErrCode != 0 {
		logger.Infof("Wechat work message sending failed: id=%v, body=%v", wechatWorkMessage.ID, string(sendMessageResponseBody))
		resendAt := time.Now().Add(5 * time.Second)
		repository.WechatWorkMessageRepo.UpdateStatusAndResendCountAndResendAtById(wechatWorkMessage.ID, repository.WECHAT_WORK_MESSAGE_STATUS_FAILED, wechatWorkMessage.ResendCount+1, &resendAt)
		repository.MessageSendingLogRepo.Save(&repository.MessageSendingLog{
			MessageID: wechatWorkMessage.ID,
			IsSuccess: false,
			Type:      repository.MESSAGE_SENDING_LOG_TYPE_WECHAT_WORK_MESSAGE,
		})
		return nil
	} else {
		logger.Infof("Wechat work message sending success. id=%v", wechatWorkMessage.ID)
		repository.WechatWorkMessageRepo.UpdateStatusById(wechatWorkMessage.ID, repository.WECHAT_WORK_MESSAGE_STATUS_SUCCESS)
		repository.MessageSendingLogRepo.Save(&repository.MessageSendingLog{
			MessageID: wechatWorkMessage.ID,
			IsSuccess: true,
			Type:      repository.MESSAGE_SENDING_LOG_TYPE_WECHAT_WORK_MESSAGE,
		})
		return nil
	}
}

type GetAccessTokenResponse struct {
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type SendWechatWorkTextMessageReq struct {
	ToUser  string                    `json:"touser"`
	MsgType string                    `json:"msgtype"`
	AgentId string                    `json:"agentid"`
	Text    SendWechatWorkTextMessage `json:"text"`
}

type SendWechatWorkTextMessage struct {
	Content string `json:"content"`
}

type SendWechatWorkTextMessageResp struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
