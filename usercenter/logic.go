package usercenter

import (
	"context"
	"errors"
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

func (uc *UserCenter) Detail(ctx context.Context, token string) (u *User, err error) {
	resp := Resp{}
	if err := uc.c.Post(uc.detailUri).SetBearerAuthToken(token).Do().Into(&resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, errors.New(resp.Msg)
	}
	detail, ok := resp.Data.(DataDetail)
	if !ok {
		return nil, errors.New("invalid data")
	}
	return &detail.User, nil
}
