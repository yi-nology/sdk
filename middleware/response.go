package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/zeromicro/go-zero/rest/httpx"
	"io"
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
		r.Response = &http.Response{

			StatusCode:       rec.statusCode,
			Proto:            r.Proto,
			ProtoMajor:       r.ProtoMajor,
			ProtoMinor:       r.ProtoMinor,
			Header:           r.Header,
			Body:             io.NopCloser(bytes.NewReader(rec.body.Bytes())),
			ContentLength:    r.ContentLength,
			TransferEncoding: r.TransferEncoding,
			Close:            r.Close,
			Trailer:          r.Trailer,
			Request:          r,
			TLS:              r.TLS,
		}
		// 将数据存入请求上下文
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
