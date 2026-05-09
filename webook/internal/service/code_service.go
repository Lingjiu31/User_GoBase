package service

import (
	"Project-WeBook/webook/internal/repository"
	"Project-WeBook/webook/internal/service/sms"
	"Project-WeBook/webook/internal/service/sms/qqemail"
	"Project-WeBook/webook/internal/service/sms/tencent"
	"context"
	"fmt"
	"math/rand"
	"strings"
)

type CodeService struct {
	tplId    string
	repo     *repository.CodeRepository
	smsSvc   *tencent.Service
	emailSvc *qqemail.Service
}

func NewCodeService(repo *repository.CodeRepository, smsSvc *tencent.Service) *CodeService {
	return &CodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *CodeService) Send(ctx context.Context, biz, target string) error {
	// 生成验证码
	code := svc.generateCode()
	// 存入 Redis
	err := svc.repo.Store(ctx, biz, target, code)
	if err != nil {
		return err
	}
	// 发送
	var sender sms.Service
	if isEmail(target) {
		// 如果是邮箱，准备邮件内容
		sender = svc.emailSvc.Ready("您的验证码", "验证码为："+code)
	} else {
		sender = svc.smsSvc.Ready(svc.tplId, []string{code})
	}

	return sender.Send(ctx, target)
}

func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *CodeService) generateCode() string {
	// 随机一个 6 位数
	num := rand.Intn(1000000)
	return fmt.Sprintf("%6d", num)
}

func isEmail(target string) bool {
	return strings.Contains(target, "@")
}
