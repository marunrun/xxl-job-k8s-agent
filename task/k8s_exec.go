package task

import (
	"context"
	"fmt"
	"github.com/marunrun/xxl-job-k8s-agent/util"
	"github.com/xxl-job/xxl-job-executor-go"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"os"
	"strconv"
)

var (
	client    *kubernetes.Clientset
	config    *rest.Config
	err       error
	namespace string
	podName   string
)

func init() {

	// 是否在集群内运行？
	inCluster, _ := strconv.ParseBool(util.GetEnv("IN_CLUSTER", "true"))

	if inCluster {
		// 使用集群内配置
		config, err = rest.InClusterConfig()
	} else {
		// 默认使用本机的~/.kube/conf 配置
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	}

	if err != nil {
		panic(err.Error())
	}

	client, err = kubernetes.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}

	namespace = os.Getenv("NAMESPACE")
	podName = os.Getenv("POD_NAME")
}

func K8s_exec(cxt context.Context, param *xxl.RunReq) string {
	logger := log.New(os.Stdout, fmt.Sprintf("XXL-AGENT-K8S-EXEC [%d]", param.LogID), 0)

	logger.Println("开始执行任务")

	// 执行命令
	cmd := []string{
		"sh",
		"-c",
		param.ExecutorParams,
	}

	// 构建请求
	req := client.CoreV1().RESTClient().Post().Resource("pods").Name(podName).
		Namespace(namespace).SubResource("exec").Param("container", os.Getenv("CONTAINER_NAME"))

	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     false,
	}

	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())

	if err != nil {
		logger.Panic(err)
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false,
	})

	if err != nil {
		logger.Panic(err)
	}

	logger.Println("任务执行完毕")
	return "executed"
}
