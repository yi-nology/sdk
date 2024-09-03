package xxl_job

type XXLJob struct {
	Enable       bool   `mapstructure:"enable,optional" json:"enable,optional" yaml:"enable"`
	ServerAddr   string `mapstructure:"serverAddr,optional" json:"serverAddr,optional" yaml:"serverAddr"`
	Token        string `json:"token,optional" yaml:"token"`
	ExecutorIp   string `json:"executorIp,optional" yaml:"executorIp"`
	ExecutorPort string `json:"executorPort,optional" yaml:"executorPort"`
	RegistryKey  string `json:"registryKey,optional" yaml:"registryKey"`
}
