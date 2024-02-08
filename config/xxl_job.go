package config

type XXLJob struct {
	ServerAddr   string `mapstructure:"server-addr" json:"server-addr" yaml:"server-addr"`
	AccessToken  string `josn:"access-token" yaml:"access-token"`
	ExecutorIp   string `json:"executor-ip" yaml:"executor-ip"`
	ExecutorPort int    `json:"executor-port" yaml:"executor-port"`
}
