package xxl_job

import (
	"github.com/go-basic/ipv4"
	"github.com/xxl-job/xxl-job-executor-go"
)

func (x XXLJob) Init() xxl.Executor {
	if x.Enable == false {
		return nil
	}
	if x.ExecutorIp == "" {
		x.ExecutorIp = ipv4.LocalIP()
	}

	exec := xxl.NewExecutor(
		xxl.ServerAddr(x.ServerAddr),     //请求地址(默认为空)
		xxl.AccessToken(x.AccessToken),   //请求令牌(默认为空)
		xxl.ExecutorPort(x.ExecutorPort), //默认9999（非必填）
		xxl.RegistryKey(x.RegistryKey),   //执行器名称
		xxl.ExecutorIp(x.ExecutorIp),     //本地IP（默认为空，可自动获取）
	)
	exec.Init()
	return exec
}
