package provider

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"time"
)

type KubernetesClient struct {
	ExternalLabels map[string]interface{}
	Cli            *kubernetes.Clientset
	Ctx            context.Context
}

func NewKubernetesClient(ctx context.Context, kubeConfigContent string, labels map[string]interface{}) (KubernetesClient, error) {
	// 如果配置内容为空，则去默认目录下取配置文件的内容
	if kubeConfigContent == "" {
		kubeConfigContent = os.Getenv("HOME") + "/.kube/config"
	}

	// 如果默认的配置文件Path实际是一个目录，那么跳过
	if _, err := os.Stat(kubeConfigContent); err == nil {
		content, err := os.ReadFile(kubeConfigContent)
		if err != nil {
			logc.Error(context.Background(), err.Error())
			return KubernetesClient{}, err
		}
		kubeConfigContent = string(content)
	}

	// 构建配置
	configBytes := []byte(kubeConfigContent)
	config, err := clientcmd.RESTConfigFromKubeConfig(configBytes)
	if err != nil {
		logc.Error(context.Background(), err.Error())
		return KubernetesClient{}, err
	}

	// 新建客户端
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		logc.Error(context.Background(), err.Error())
	}

	return KubernetesClient{
		Cli:            cs,
		Ctx:            ctx,
		ExternalLabels: labels,
	}, nil
}

func (a KubernetesClient) GetWarningEvent(reason string, scope int) (*corev1.EventList, error) {
	var warningEvents corev1.EventList
	cutoffTime := time.Now().Add(-time.Duration(scope) * time.Minute)
	opts := metav1.ListOptions{
		Limit:         50, // 减少每次请求的数量，防止过多资源占用
		FieldSelector: "reason=" + reason,
	}

	for {
		list, err := a.Cli.CoreV1().Events(corev1.NamespaceAll).List(a.Ctx, opts)
		if err != nil {
			return nil, err
		}

		for _, event := range list.Items {
			// 检查事件的 Reason 和事件发生时间
			eventTime := event.EventTime
			if event.Reason == reason && eventTime.After(cutoffTime) {
				warningEvents.Items = append(warningEvents.Items, event)
			}
		}

		// 如果没有更多事件，则停止拉取
		if list.Continue == "" {
			break
		}

		// 使用 Continue 获取下一页
		opts.Continue = list.Continue
	}

	return &warningEvents, nil
}

func (a KubernetesClient) GetExternalLabels() map[string]interface{} {
	return a.ExternalLabels
}

func (a KubernetesClient) Check() (bool, error) {
	_, err := a.Cli.ServerVersion()

	return err == nil, err
}
