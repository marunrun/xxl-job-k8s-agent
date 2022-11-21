package main

import (
	"fmt"
	"github.com/marunrun/xxl-job-k8s-agent/task"
	"github.com/marunrun/xxl-job-k8s-agent/util"
	"github.com/xxl-job/xxl-job-executor-go"
	"log"
)

var (
	serverAddr  string
	accessToken string
	jobName     string
)

func init() {
	serverAddr = util.GetEnv("XXL_JOB_ADDR", "")
	accessToken = util.GetEnv("XXL_JOB_ACCESS_TOKEN", "")
	jobName = util.GetEnv("XXL_JOB_NAME", "")
	if len(jobName) < 1 {
		panic("XXL_JOB_NAME is empty")
	}
}

func main() {
	exec := xxl.NewExecutor(
		xxl.ServerAddr(serverAddr),
		xxl.AccessToken(accessToken), //请求令牌(默认为空)
		xxl.ExecutorPort("9999"),     //默认9999（非必填）
		xxl.RegistryKey(jobName),     //执行器名称
		xxl.SetLogger(&logger{}),     //自定义日志
	)
	exec.Init()

	//设置日志查看handler
	exec.LogHandler(func(req *xxl.LogReq) *xxl.LogRes {
		return &xxl.LogRes{Code: xxl.SuccessCode, Msg: "", Content: xxl.LogResContent{
			FromLineNum: req.FromLineNum,
			ToLineNum:   2,
			LogContent:  fmt.Sprintf("请根据logId [%d] 到阿里云日志中心查看日志", req.LogID),
			IsEnd:       true,
		}}
	})

	//注册任务handler
	exec.RegTask("k8s.exec", task.K8s_exec)
	log.Fatal(exec.Run())
}

//xxl.Logger接口实现
type logger struct{}

func (l *logger) Info(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf("XXL-AGENT  "+format, a...))
}

func (l *logger) Error(format string, a ...interface{}) {
	log.Println(fmt.Sprintf("XXL-AGENT  "+format, a...))
}
