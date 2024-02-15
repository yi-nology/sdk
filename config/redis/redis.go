package redis

type Redis struct {
	Enable   bool   `mapstructure:"enable" json:"enable,optional" yaml:"enable"`       // 是否开启
	Addr     string `mapstructure:"addr" json:"addr,optional" yaml:"addr"`             // 服务器地址:端口
	Password string `mapstructure:"password" json:"password,optional" yaml:"password"` // 密码
	DB       int    `mapstructure:"db" json:"db,optional" yaml:"db"`                   // redis的哪个数据库
}
