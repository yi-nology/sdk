package client

import (
	"context"
	"github.com/apolloconfig/agollo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/xxl-job/xxl-job-executor-go"
	"github.com/yi-nology/sdk/conf/apollo"
	"github.com/yi-nology/sdk/conf/mongo"
	"github.com/yi-nology/sdk/conf/mysql"
	redisConf "github.com/yi-nology/sdk/conf/redis"
	"github.com/yi-nology/sdk/conf/xxl_job"
	"gorm.io/gorm"
)

var (
	// ApolloConfig apollo配置
	Apollo agollo.Client
	// MysqlConfig mysql配置
	Mysql *gorm.DB
	// RedisConfig redis配置
	Redis *redis.Client
	// XXLJobConfig xxl-job配置
	XXLJob xxl.Executor
	// MongoClient mongo配置
	Mongo *mongo.MongoCli

	GloablConfig *Config
)

type Config struct {
	Apollo apollo.Apollo   `mapstructure:"apollo,optional" json:"apollo,optional" yaml:"apollo,optional"`
	Mysql  mysql.Mysql     `mapstructure:"mysql,optional" json:"mysql,optional" yaml:"mysql,optional"`
	Redis  redisConf.Redis `mapstructure:"redis,optional" json:"redis,optional" yaml:"redis,optional"`
	XXLJob xxl_job.XXLJob  `mapstructure:"xxlJob,optional" json:"xxlJob,optional" yaml:"xxlJob,optional"`
	Mongo  mongo.Mongo     `mapstructure:"mongo,optional" json:"mongo,optional" yaml:"mongo,optional"`
}

func (i *Config) Init() (err error) {
	Apollo, err = i.Apollo.InitApollo()
	if err != nil {
		return err
	}
	Mysql, err = i.Mysql.Init()
	if err != nil {
		return err
	}
	Redis, err = i.Redis.Init()
	if err != nil {
		return err
	}
	Mongo, err = i.Mongo.Init()
	if err != nil {
		return err
	}
	XXLJob = i.XXLJob.Init()

	GloablConfig = i
	return nil
}

func (i *Config) Stop(ctx context.Context) {
	if i.Mongo.Enable {
		Mongo.Close(ctx)
	}
	if i.Mysql.Enable {
	}
	if i.Redis.Enable {
		Redis.Close()
	}
	if i.Apollo.Enable {
		Apollo.Close()
	}
	if i.XXLJob.Enable {
		XXLJob.Stop()
	}
}
