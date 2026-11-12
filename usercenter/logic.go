package usercenter

import (
	"context"
	"errors"

	"github.com/imroc/req/v3"
)

type UserCenter struct {
	c         *req.Client
	detailUri string
}

// NewUserCenter 按传入参数构造一个新的 UserCenter 实例。
// 每次调用都返回独立实例（可同时连接多个用户中心环境），
// 不再做包级单例缓存：此前用包级 sync.Once 导致第二次构造的
// baseUrl/detailUri 被静默丢弃，请求会打到第一个实例的服务上。
func NewUserCenter(baseUrl, detailUri string) *UserCenter {
	return &UserCenter{
		c:         req.C().SetBaseURL(baseUrl),
		detailUri: detailUri,
	}
}

func (uc *UserCenter) Detail(ctx context.Context, token string) (u *User, err error) {
	resp := Resp{}
	// ctx 必须透传给请求，否则调用方的取消/超时完全失效
	if err := uc.c.Post(uc.detailUri).SetContext(ctx).SetBearerAuthToken(token).Do().Into(&resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, errors.New(resp.Msg)
	}
	return &resp.Data.User, nil
}
