package sms

import "context"

type Service interface {
	Send(ctx context.Context, target ...string) error
}
