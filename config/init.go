package config

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/yi-nology/sdk/config/apollo"
	"github.com/yi-nology/sdk/config/mysql"
	redisConf "github.com/yi-nology/sdk/config/redis"
	"gorm.io/gorm"
)

var (
	// ApolloConfig apollo配置
	ApolloConfig *agollo.Client
	// MysqlConfig mysql配置
	MysqlConfig *gorm.DB
	// RedisConfig redis配置
	RedisClient *redis.Client
)

type Config struct {
	Apollo apollo.Apollo   `mapstructure:"apollo" json:"apollo" yaml:"apollo"`
	Mysql  mysql.Mysql     `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis  redisConf.Redis `mapstructure:"redis" json:"redis" yaml:"redis"`
	XXLJob XXLJob          `mapstructure:"xxl-job" json:"xxl-job" yaml:"xxl-job"`
}

func (i *Config) Init() (err error) {
	ApolloConfig, err = i.Apollo.InitApollo()
	if err != nil {
		return err
	}
	MysqlConfig, err = i.Mysql.Init()
	if err != nil {
		return err
	}
	RedisClient, err = i.Redis.Init()
	if err != nil {
		return err
	}
	return nil
}
