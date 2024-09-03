package xxl_job

import (
	"github.com/xxl-job/xxl-job-executor-go"
)

func (x XXLJob) Init() xxl.Executor {
	if x.Enable == false {
		return nil
	}
	option := []xxl.Option{xxl.ServerAddr(x.ServerAddr), xxl.RegistryKey(x.RegistryKey)}
	if x.ExecutorIp != "" {
		option = append(option, xxl.ExecutorIp(x.ExecutorIp))
	}
	if x.Token != "" {
		option = append(option, xxl.AccessToken(x.Token))
	}
	if x.ExecutorPort != "" {
		option = append(option, xxl.ExecutorPort(x.ExecutorPort))
	}

	exec := xxl.NewExecutor(option...)
	exec.Init()
	return exec
}
