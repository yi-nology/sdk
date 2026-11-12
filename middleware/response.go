package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type responseWrapper struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ResponseWithLogMiddleware struct {
}

func NewResponseMiddleware() *ResponseWithLogMiddleware {
	return &ResponseWithLogMiddleware{}
}

// debugHeaderWhitelist 是访问日志允许落盘的低敏请求头白名单。
// 采用白名单而非黑名单：凭据类头（Authorization/Cookie/X-Token 等）及任何
// 未收录的头一律不落盘，避免新增敏感头时日志泄漏。
var debugHeaderWhitelist = []string{
	"Accept",
	"Content-Length",
	"Content-Type",
	"User-Agent",
	"X-Forwarded-For",
	"X-Request-Id",
}

// maxLoggedBodyBytes 限制日志里响应体的最大长度，避免大响应撑爆日志量。
const maxLoggedBodyBytes = 2 * 1024

// debugHeaders 按白名单过滤请求头，只保留低敏字段用于排障。
func debugHeaders(h http.Header) map[string][]string {
	out := make(map[string][]string, len(h))
	for k, v := range h {
		if slices.ContainsFunc(debugHeaderWhitelist, func(allow string) bool {
			return strings.EqualFold(k, allow)
		}) {
			out[k] = v
		}
	}
	return out
}

// truncateBody 截断超长响应体，并标注被截断的字节数。
func truncateBody(s string) string {
	if len(s) <= maxLoggedBodyBytes {
		return s
	}
	return fmt.Sprintf("%s...(%d bytes truncated)", s[:maxLoggedBodyBytes], len(s)-maxLoggedBodyBytes)
}

func (m *ResponseWithLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 创建一个新的响应记录器
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next(rec, r)

		// 请求头走低敏白名单，响应体做截断；r.Form/r.Body 可能携带凭据，
		// 且 r.Body 打印出来只是结构体指针，一并从日志移除。
		accessLog := fmt.Sprintf("Method：%s, Url：%s, Header：%+v; code:%+v, data:%s",
			r.Method, r.URL, debugHeaders(r.Header), rec.statusCode, truncateBody(rec.body.String()))

		// 访问日志降到 Debug 级：常规级别下不再为每个请求落一份头/体副本
		logx.WithContext(r.Context()).Debugf("%s", accessLog)

		// 检查响应状态码，如果是错误码则返回失败信息
		if rec.statusCode >= 400 {
			logx.WithContext(r.Context()).Errorf("%s", accessLog)
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
