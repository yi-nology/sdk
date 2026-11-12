package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/zeromicro/go-zero/core/logx"
)

// captureWriter 实现 go-zero logx.Writer，把日志捕获进内存供断言。
type captureWriter struct {
	mu   sync.Mutex
	logs bytes.Buffer
}

func (w *captureWriter) write(v any) {
	w.mu.Lock()
	defer w.mu.Unlock()
	fmt.Fprint(&w.logs, v)
}

func (w *captureWriter) Alert(v any)                     { w.write(v) }
func (w *captureWriter) Close() error                    { return nil }
func (w *captureWriter) Debug(v any, _ ...logx.LogField) { w.write(v) }
func (w *captureWriter) Error(v any, _ ...logx.LogField) { w.write(v) }
func (w *captureWriter) Info(v any, _ ...logx.LogField)  { w.write(v) }
func (w *captureWriter) Severe(v any)                    { w.write(v) }
func (w *captureWriter) Slow(v any, _ ...logx.LogField)  { w.write(v) }
func (w *captureWriter) Stack(v any)                     { w.write(v) }
func (w *captureWriter) Stat(v any, _ ...logx.LogField)  { w.write(v) }

func (w *captureWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.logs.String()
}

func setupCapture(t *testing.T, level uint32) *captureWriter {
	t.Helper()
	cw := &captureWriter{}
	logx.SetWriter(cw)
	logx.SetLevel(level)
	t.Cleanup(func() {
		logx.Reset()
		logx.SetLevel(logx.InfoLevel)
	})
	return cw
}

const (
	testAuth   = "Bearer secret-token-abc"
	testCookie = "session=cookie-secret"
	testRespTk = "resp-token-secret"
)

func doLoggedRequest(t *testing.T, status int, body string) {
	t.Helper()
	mw := NewResponseMiddleware()
	handler := mw.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.Header.Set("Authorization", testAuth)
	req.Header.Set("Cookie", testCookie)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", "req-42")
	handler(rec, req)
}

func assertNoLeak(t *testing.T, logs string) {
	t.Helper()
	for _, secret := range []string{testAuth, testCookie, testRespTk, "secret-token"} {
		if strings.Contains(logs, secret) {
			t.Errorf("日志泄漏凭据 %q", secret)
		}
	}
}

// 回归 #4：访问日志降到 Debug 级后，Info 级别下不再为每个请求落盘头与体。
func TestAccessLogDowngradedToDebug(t *testing.T) {
	cw := setupCapture(t, logx.InfoLevel)
	doLoggedRequest(t, http.StatusOK, `{"token":"`+testRespTk+`"}`)
	if got := cw.String(); got != "" {
		t.Fatalf("Info 级别下访问日志不应落盘, got: %s", got)
	}
}

// 回归 #4：Debug 级别下日志可用但凭据必须脱敏，白名单低敏头保留，长响应体截断。
// 注：Debug 为显式开启的排障通道，响应体按 issue 口径做截断（含凭据的接口
// 应自行按接口维度关闭 body 日志）；默认 Info 级别下头与体完全不落盘。
func TestAccessLogSanitizedAtDebugLevel(t *testing.T) {
	cw := setupCapture(t, logx.DebugLevel)
	doLoggedRequest(t, http.StatusOK, `{"data":"`+strings.Repeat("x", 8*1024)+`"}`)

	logs := cw.String()
	assertNoLeak(t, logs)

	if !strings.Contains(logs, "Method：POST") {
		t.Error("Debug 级别下应保留访问日志行")
	}
	for _, allow := range []string{"Content-Type", "X-Request-Id"} {
		if !strings.Contains(logs, allow) {
			t.Errorf("白名单低敏头 %s 应保留用于排障", allow)
		}
	}
	if strings.Contains(logs, "X-TOKEN") {
		t.Error("未收录的请求头不应出现在日志里")
	}
	if !strings.Contains(logs, "bytes truncated") {
		t.Error("超长响应体应被截断并标注")
	}
}

// 回归 #4：>=400 的 Error 分支同样不得泄漏凭据。
func TestErrorBranchSanitized(t *testing.T) {
	cw := setupCapture(t, logx.InfoLevel)
	doLoggedRequest(t, http.StatusInternalServerError, `{"error":"boom"}`)

	logs := cw.String()
	assertNoLeak(t, logs)

	if !strings.Contains(logs, "Method：POST") {
		t.Error("错误分支应保留日志行用于排障")
	}
	if !strings.Contains(logs, "Content-Type") {
		t.Error("错误分支应保留白名单低敏头")
	}
}
