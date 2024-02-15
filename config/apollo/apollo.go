package apollo

type Apollo struct {
	Enable bool   `json:"enable,optional" yaml:"enable,optional" default:"false"`
	AppID  string `json:"appId,optional" yaml:"appId,optional" default:""`
	// Cluster is the cluster name, default is "default".
	Cluster string `json:"cluster,optional" yaml:"cluster,optional" default:"default"`
	// NamespaceName is the namespace name, default is "application".
	NamespaceName  string `json:"namespaceName,optional" yaml:"namespaceName,optional" default:"application"`
	Endpoint       string `json:"endpoint,optional" yaml:"endpoint,optional" default:"http://localhost:8080"`
	Dynamic        bool   `json:"dynamic,optional" yaml:"dynamic,optional" yaml:"dynamic"`
	Interval       int    `json:"interval,optional" yaml:"interval,optional" default:"30s"`
	Secret         string `json:"secret,optional" yaml:"secret,optional"`
	IsBackupConfig bool   `json:"isBackupConfig,optional" yaml:"isBackupConfig,optional" default:"true"`
}
