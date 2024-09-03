package mongo

import (
	"context"
	"fmt"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/event"
	opt "go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"strings"
)

type Mongo struct {
	Enable           bool             `json:"enable,optional" yaml:"enable,optional" default:"false"`
	Coll             string           `json:"coll" yaml:"coll" mapstructure:"coll"`                                           // collection name
	Options          string           `json:"options" yaml:"options" mapstructure:"options"`                                  // mongodb options
	Database         string           `json:"database" yaml:"database" mapstructure:"database"`                               // database name
	Username         string           `json:"username" yaml:"username" mapstructure:"username"`                               // 用户名
	Password         string           `json:"password" yaml:"password" mapstructure:"password"`                               // 密码
	AuthSource       string           `json:"auth-source" yaml:"auth-source" mapstructure:"auth-source"`                      // 验证数据库
	MinPoolSize      uint64           `json:"min-pool-size" yaml:"min-pool-size" mapstructure:"min-pool-size"`                // 最小连接池
	MaxPoolSize      uint64           `json:"max-pool-size" yaml:"max-pool-size" mapstructure:"max-pool-size"`                // 最大连接池
	SocketTimeoutMs  int64            `json:"socket-timeout-ms" yaml:"socket-timeout-ms" mapstructure:"socket-timeout-ms"`    // socket超时时间
	ConnectTimeoutMs int64            `json:"connect-timeout-ms" yaml:"connect-timeout-ms" mapstructure:"connect-timeout-ms"` // 连接超时时间
	IsZap            bool             `json:"is-zap" yaml:"is-zap" mapstructure:"is-zap"`                                     // 是否开启zap日志
	Hosts            []*MongoHost     `json:"hosts" yaml:"hosts" mapstructure:"hosts"`                                        // 主机列表
	MongoClient      *qmgo.QmgoClient `json:"-,optional" yaml:"-,optional" mapstructure:"-,optional"`
}

type MongoHost struct {
	Host string `json:"host" yaml:"host" mapstructure:"host"` // ip地址
	Port string `json:"port" yaml:"port" mapstructure:"port"` // 端口
}

// Uri .
func (m *Mongo) Uri() string {
	length := len(m.Hosts)
	hosts := make([]string, 0, length)
	for i := 0; i < length; i++ {
		if m.Hosts[i].Host != "" && m.Hosts[i].Port != "" {
			hosts = append(hosts, m.Hosts[i].Host+":"+m.Hosts[i].Port)
		}
	}
	if m.Options != "" {
		return fmt.Sprintf("mongodb://%s/%s?%s", strings.Join(hosts, ","), m.Database, m.Options)
	}
	return fmt.Sprintf("mongodb://%s/%s", strings.Join(hosts, ","), m.Database)
}

func (m *Mongo) GetClientOptions() []options.ClientOptions {
	cmdMonitor := &event.CommandMonitor{
		Started: func(ctx context.Context, event *event.CommandStartedEvent) {
			zap.L().Info(fmt.Sprintf("[MongoDB][RequestID:%d][database:%s] %s\n", event.RequestID, event.DatabaseName, event.Command), zap.String("business", "mongo"))
		},
		Succeeded: func(ctx context.Context, event *event.CommandSucceededEvent) {
			zap.L().Info(fmt.Sprintf("[MongoDB][RequestID:%d] [%s] %s\n", event.RequestID, event.Duration.String(), event.Reply), zap.String("business", "mongo"))
		},
		Failed: func(ctx context.Context, event *event.CommandFailedEvent) {
			zap.L().Error(fmt.Sprintf("[MongoDB][RequestID:%d] [%s] %s\n", event.RequestID, event.Duration.String(), event.Failure), zap.String("business", "mongo"))
		},
	}
	return []options.ClientOptions{{ClientOptions: &opt.ClientOptions{Monitor: cmdMonitor}}}
}
