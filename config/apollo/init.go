package apollo

import (
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
)

func (a Apollo) InitApollo() (*agollo.Client, error) {

	var c = &config.AppConfig{
		AppID:          a.AppID,
		Cluster:        a.Cluster,
		IP:             a.Endpoint,
		NamespaceName:  a.NamespaceName,
		IsBackupConfig: true,
		Secret:         a.Secret,
	}
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("初始化Apollo配置成功")
	return &client, err
}
