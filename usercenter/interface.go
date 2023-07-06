package usercenter

import "context"

type UserCenterApi interface {
	Detail(ctx context.Context, token string) (u *User, err error)
}
