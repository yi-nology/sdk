package config

type XXLJob struct {
	ServerAddr   string `mapstructure:"serverAddr,optional" json:"serverAddr,optional" yaml:"serverAddr"`
	AccessToken  string `josn:"accessToken,optional" yaml:"accessToken"`
	ExecutorIp   string `json:"executorIp,optional" yaml:"executorIp"`
	ExecutorPort int    `json:"executorPort,optional" yaml:"executorPort"`
}
