package middleware

import (
	"bytes"
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

		// 检查响应状态码，如果是错误码则返回失败信息
		if rec.statusCode >= 400 {
			httpx.OkJson(w, responseWrapper{
				Code: 1,
				Msg:  rec.body.String(),
				Data: nil,
			})
			return
		}
		// 封装原始响应数据
		var originalData interface{}
		if err := json.Unmarshal(rec.body.Bytes(), &originalData); err != nil {
			httpx.OkJson(w, responseWrapper{
				Code: 1,
				Msg:  "失败",
				Data: nil,
			})
			return
		}

		// 返回封装后的响应数据
		httpx.OkJson(w, responseWrapper{
			Code: 0,
			Msg:  "ok",
			Data: originalData,
		})
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
