package mysql

type Mysql struct {
	GeneralDB `yaml:",inline" mapstructure:",squash"`
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}

func (m *Mysql) GetLogMode() string {
	return m.LogMode
}

type DsnProvider interface {
	Dsn() string
}

type GeneralDB struct {
	Enable       bool   `mapstructure:"enable" json:"enable,optional" yaml:"enable"`
	Prefix       string `mapstructure:"prefix" json:"prefix,optional" yaml:"prefix"`
	Port         string `mapstructure:"port" json:"port,optional" yaml:"port"`
	Config       string `mapstructure:"config" json:"config,optional" yaml:"config"`       // 高级配置
	Dbname       string `mapstructure:"db-name" json:"db-name,optional" yaml:"db-name"`    // 数据库名
	Username     string `mapstructure:"username" json:"username,optional" yaml:"username"` // 数据库密码
	Password     string `mapstructure:"password" json:"password,optional" yaml:"password"` // 数据库密码
	Path         string `mapstructure:"path" json:"path,optional" yaml:"path"`
	Engine       string `mapstructure:"engine" json:"engine,optional" yaml:"engine" default:"InnoDB"`        //数据库引擎，默认InnoDB
	LogMode      string `mapstructure:"log-mode" json:"log-mode,optional" yaml:"log-mode"`                   // 是否开启Gorm全局日志
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns,optional" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns,optional" yaml:"max-open-conns"` // 打开到数据库的最大连接数
	Singular     bool   `mapstructure:"singular" json:"singular,optional" yaml:"singular"`                   //是否开启全局禁用复数，true表示开启
	LogZap       bool   `mapstructure:"log-zap" json:"log-zap,optional" yaml:"log-zap"`                      // 是否通过zap写入日志文件
}
