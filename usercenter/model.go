package usercenter

type User struct {
	Id       int64  `json:"id"`
	Mobile   string `json:"mobile"`
	Nickname string `json:"nickname"`
	Sex      int64  `json:"sex"`
	Avatar   string `json:"avatar"`
	Info     string `json:"info"`
}

// Resp 是用户中心的标准响应包裹。
// Data 声明为具体类型 DataDetail，让 encoding/json 直接解码；
// 此前声明为 interface{} 时，JSON 解出来必然是 map[string]interface{}，
// 对它做 DataDetail 类型断言恒失败，合法响应也会返回 "invalid data"。
type Resp struct {
	Code int64      `json:"code"`
	Msg  string     `json:"msg"`
	Data DataDetail `json:"data"`
}

type DataDetail struct {
	User User `json:"user"`
}
