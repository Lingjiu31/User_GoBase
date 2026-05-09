package tencent

import (
	"context"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId    string
	signName string
	smsSend  *smsSend
	client   *sms.Client
}

type smsSend struct {
	tpl  string
	args []string
}

func NewService(appId string, signName string, client *sms.Client) *Service {
	return &Service{
		appId:    appId,
		signName: signName,
		client:   client,
	}
}

func (s *Service) Ready(tpl string, args []string) *Service {
	newService := *s
	newService.smsSend = &smsSend{
		tpl:  tpl,
		args: args,
	}
	return &newService
}

func (s *Service) Send(ctx context.Context, target ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = common.StringPtr(s.appId)
	req.SignName = common.StringPtr(s.signName)
	req.TemplateId = common.StringPtr(s.smsSend.tpl)
	req.TemplateParamSet = common.StringPtrs(s.smsSend.args)
	req.PhoneNumberSet = common.StringPtrs(target)

	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "OK" {
			return fmt.Errorf("短信发送失败 %s, %s", *status.Code, *status.Message)
		}
	}
	return nil
}
