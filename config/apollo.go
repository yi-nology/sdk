package config

type Apollo struct {
	Enable bool   `json:"enable" yaml:"enable" default:"false"`
	AppID  string `json:"appId" yaml:"appId" default:""`
	// Cluster is the cluster name, default is "default".
	Cluster string `json:"cluster" yaml:"cluster" default:"default"`
	// NamespaceName is the namespace name, default is "application".
	NamespaceName string `json:"namespaceName" yaml:"namespaceName" default:"application"`
	Endpoint      string `json:"endpoint" yaml:"endpoint" default:"http://localhost:8080"`
	Dynamic       bool   `json:"dynamic" yaml:"dynamic" yaml:"dynamic"`
	Interval      int    `json:"interval" yaml:"interval" default:"30s"`
	Secret        string `json:"secret" yaml:"secret"`
}
