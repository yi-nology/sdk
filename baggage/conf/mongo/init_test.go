package mongo

import (
	"testing"

	"github.com/qiniu/qmgo"
)

func TestParseIndexKeys(t *testing.T) {
	cases := []struct {
		name     string
		keys     []string
		wantSpec []string
		wantName string
	}{
		{
			name:     "单键升序",
			keys:     []string{"user_id"},
			wantSpec: []string{"user_id"},
			wantName: "user_id_1",
		},
		{
			name:     "多键混合方向",
			keys:     []string{"-created_at", "name"},
			wantSpec: []string{"-created_at", "name"},
			wantName: "created_at_-1_name_1",
		},
		{
			name:     "多键全升序",
			keys:     []string{"a", "b"},
			wantSpec: []string{"a", "b"},
			wantName: "a_1_b_1",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			spec, name, err := parseIndexKeys(tc.keys)
			if err != nil {
				t.Fatalf("不应返回错误, got: %v", err)
			}
			if name != tc.wantName {
				t.Errorf("规范名 got %q, want %q", name, tc.wantName)
			}
			if len(spec) != len(tc.wantSpec) {
				t.Fatalf("spec got %v, want %v", spec, tc.wantSpec)
			}
			for i := range spec {
				if spec[i] != tc.wantSpec[i] {
					t.Errorf("spec[%d] got %q, want %q", i, spec[i], tc.wantSpec[i])
				}
			}
		})
	}
}

// 回归 #5：两个单键索引的配置在旧实现中必 panic
// （indexes[i][i] 用外层下标访问内层切片，i=1 时 indexes[1][1] 越界）。
func TestParseIndexKeysMultipleSingleKeySpecsDoNotPanic(t *testing.T) {
	indexes := [][]string{{"user_id"}, {"created_at"}}
	wantNames := []string{"user_id_1", "created_at_1"}
	for i := range indexes {
		spec, name, err := parseIndexKeys(indexes[i])
		if err != nil {
			t.Fatalf("第%d个索引定义不应返回错误, got: %v", i, err)
		}
		if name != wantNames[i] {
			t.Errorf("第%d个索引规范名 got %q, want %q", i, name, wantNames[i])
		}
		if len(spec) != 1 || spec[0] != indexes[i][0] {
			t.Errorf("第%d个索引 spec got %v, want %v", i, spec, indexes[i])
		}
	}
}

// 方向必须真正传给驱动：spec 保留 qmgo 的 "-" 前缀约定，
// qmgo.SplitSortField（ensureIndex 的同款解析）应得出正确的字段名与方向。
func TestParseIndexKeysDirectionReachesDriverConvention(t *testing.T) {
	spec, _, err := parseIndexKeys([]string{"-created_at", "name"})
	if err != nil {
		t.Fatalf("不应返回错误, got: %v", err)
	}
	want := map[string]int32{"created_at": -1, "name": 1}
	for _, k := range spec {
		key, sort := qmgo.SplitSortField(k)
		if want[key] != sort {
			t.Errorf("键 %q 解出 %q:%d, want %d", k, key, sort, want[key])
		}
	}
}

func TestParseIndexKeysInvalidConfig(t *testing.T) {
	cases := []struct {
		name string
		keys []string
	}{
		{name: "空索引定义", keys: []string{}},
		{name: "空字符串键", keys: []string{"a", ""}},
		{name: "只有降序标记", keys: []string{"-"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := parseIndexKeys(tc.keys); err == nil {
				t.Error("非法配置应返回明确错误而不是 panic")
			}
		})
	}
}
