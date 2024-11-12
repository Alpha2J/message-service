package manager

import (
	"message-service/internal/app/repository"
	"message-service/internal/pkg/helper"
	"strconv"
)

type WechatWorkAppConfigManager struct{}

var WechatWorkAppConfigMana WechatWorkAppConfigManager = WechatWorkAppConfigManager{}

var wechatWorkAppConfigCache *helper.Cache = helper.NewCache()

func (*WechatWorkAppConfigManager) FindById(id int64) (*repository.WechatWorkAppConfig, error) {
	value, found := wechatWorkAppConfigCache.Get(strconv.FormatInt(id, 10))
	if found {
		return value.(*repository.WechatWorkAppConfig), nil
	}

	wechatWorkAppConfig, err := repository.WechatWorkAppConfigRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	wechatWorkAppConfigCache.Set(strconv.FormatInt(id, 10), wechatWorkAppConfig)

	return wechatWorkAppConfig, nil
}
