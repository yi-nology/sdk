package usercenter

import (
	"context"
	"github.com/imroc/req/v3"
	"sync"
)

var once sync.Once
var uc *UserCenter

type UserCenter struct {
	c         *req.Client
	detailUri string
}

func NewUserCenter(baseUrl, detailUri string) *UserCenter {
	once.Do(func() {
		uc = &UserCenter{
			c:         req.C().SetBaseURL(baseUrl),
			detailUri: detailUri,
		}
	})
	return uc
}

func (uc *UserCenter) Detail(ctx context.Context, token string) (u *Resp, err error) {
	if err := uc.c.Post(uc.detailUri).SetBearerAuthToken(token).Do().Into(&u); err != nil {
		return nil, err
	}
	return u, nil
}
