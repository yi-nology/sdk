package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

type responseWrapper struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ResponseMiddleware struct {
}

func NewResponseMiddleware() *ResponseMiddleware {
	return &ResponseMiddleware{}
}

func (m *ResponseMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 创建一个新的响应记录器
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next(rec, r)

		// 将数据存入请求上下文
		rec.setCtx(r)
		// 检查响应状态码，如果是错误码则返回失败信息
		if rec.statusCode >= 400 {
			httpx.OkJson(w, responseWrapper{
				Code: 1,
				Msg:  rec.body.String(),
				Data: nil,
			})
			return
		}

		if rec.body.Bytes() == nil {
			httpx.OkJson(w, responseWrapper{
				Code: 0,
				Msg:  "成功",
				Data: nil,
			})
			return
		}
		// 封装原始响应数据
		wrapper := responseWrapper{
			Code: 0,
			Msg:  "成功",
		}
		if err := json.Unmarshal(rec.body.Bytes(), &wrapper.Data); err != nil {
			wrapper.Code = 1
			wrapper.Msg = "数据解析失败"
		}
		// 返回封装后的响应数据
		httpx.OkJson(w, wrapper)
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
}

func (rec *responseRecorder) Write(body []byte) (int, error) {
	rec.body.Write(body)
	return len(body), nil
}

func (rec *responseRecorder) setCtx(r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), "responseInfo", rec))
}

func GetResponse(r *http.Request) (int, bytes.Buffer) {
	info, ok := r.Context().Value("responseInfo").(*responseRecorder)
	if !ok {
		return 0, bytes.Buffer{}
	}
	return info.statusCode, info.body
}
