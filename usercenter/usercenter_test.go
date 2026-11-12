package usercenter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"
)

// newDetailServer 起一个返回固定响应体的用户中心假服务，并统计收到的请求数。
func newDetailServer(t *testing.T, body string, hits *atomic.Int64) *httptest.Server {
	t.Helper()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hits != nil {
			hits.Add(1)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(s.Close)
	return s
}

func detailBody(id int64, nickname string) string {
	return `{"code":0,"msg":"ok","data":{"user":{"id":` +
		strconv.FormatInt(id, 10) + `,"mobile":"13800000000","nickname":"` + nickname + `"}}}`
}

// 回归 #1：合法响应（code=0、data.user 齐全）必须解出 User，
// 旧实现 Resp.Data 为 interface{}，对 map[string]interface{} 做
// DataDetail 类型断言恒失败，Detail 100% 返回 "invalid data"。
func TestDetailParsesValidResponse(t *testing.T) {
	s := newDetailServer(t, detailBody(42, "alice"), nil)
	uc := NewUserCenter(s.URL, "/detail")

	u, err := uc.Detail(context.Background(), "tok")
	if err != nil {
		t.Fatalf("合法响应不应返回错误, got: %v", err)
	}
	if u == nil {
		t.Fatal("应返回非 nil 的 User")
	}
	if u.Id != 42 || u.Nickname != "alice" {
		t.Fatalf("解码结果不对: got %+v, want id=42 nickname=alice", u)
	}
}

// 回归 #2：第二次构造的 baseUrl/detailUri 不允许被包级 sync.Once 静默丢弃。
func TestNewUserCenterReturnsIndependentInstances(t *testing.T) {
	sa := newDetailServer(t, detailBody(1, "from-a"), nil)
	sb := newDetailServer(t, detailBody(2, "from-b"), nil)

	ucA := NewUserCenter(sa.URL, "/detail")
	ucB := NewUserCenter(sb.URL, "/detail")

	if ucA == ucB {
		t.Fatal("两次构造应返回不同实例，而不是同一个包级单例")
	}

	ub, err := ucB.Detail(context.Background(), "tok")
	if err != nil {
		t.Fatalf("第二次构造的实例请求失败: %v", err)
	}
	if ub.Nickname != "from-b" {
		t.Fatalf("第二次构造的 baseUrl 未生效: got nickname=%q, want from-b", ub.Nickname)
	}

	ua, err := ucA.Detail(context.Background(), "tok")
	if err != nil {
		t.Fatalf("第一次构造的实例请求失败: %v", err)
	}
	if ua.Nickname != "from-a" {
		t.Fatalf("第一次构造的 baseUrl 未生效: got nickname=%q, want from-a", ua.Nickname)
	}
}

// 回归 #3：已取消的 ctx 必须让请求发不出去（服务端收到的请求数为 0）。
func TestDetailHonorsContextCancellation(t *testing.T) {
	var hits atomic.Int64
	s := newDetailServer(t, detailBody(1, "alice"), &hits)
	uc := NewUserCenter(s.URL, "/detail")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := uc.Detail(ctx, "tok")
	if err == nil {
		t.Fatal("ctx 已取消时 Detail 应返回错误")
	}
	if n := hits.Load(); n != 0 {
		t.Fatalf("ctx 已取消仍打出了 %d 次请求, want 0", n)
	}
}

// 业务错误码透传：code!=0 时应返回服务端的 Msg。
func TestDetailPropagatesErrorCode(t *testing.T) {
	s := newDetailServer(t, `{"code":401,"msg":"token expired"}`, nil)
	uc := NewUserCenter(s.URL, "/detail")

	u, err := uc.Detail(context.Background(), "tok")
	if err == nil {
		t.Fatal("code!=0 应返回错误")
	}
	if err.Error() != "token expired" {
		t.Fatalf("错误信息应为服务端 Msg, got %q", err.Error())
	}
	if u != nil {
		t.Fatal("错误时不应返回 User")
	}
}
