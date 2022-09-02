package services

import (
	"context"
	"strings"

	"github.com/cjlapao/postfixcli-backend-api/ioc"
	"github.com/cjlapao/postfixcli-backend-api/models"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var globalK8sService *KubernetesClusterService

type KubernetesClusterService struct {
	Context            context.Context
	KubeConfigPath     string
	UseInClusterConfig bool
	Client             *kubernetes.Clientset
}

func GetK8sService() *KubernetesClusterService {
	if globalK8sService != nil {
		return globalK8sService
	}

	return NewK8sService()
}

func NewK8sService() *KubernetesClusterService {
	globalK8sService = &KubernetesClusterService{}

	globalK8sService.Context = context.Background()

	kubeConfig := ioc.Config.GetString("kubeconfig")
	if kubeConfig != "" {
		globalK8sService.KubeConfigPath = kubeConfig
	}

	globalK8sService.initClient()
	return globalK8sService
}

func (svc *KubernetesClusterService) initClient() {
	var config *rest.Config
	var err error

	if svc.KubeConfigPath == "" {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", svc.KubeConfigPath)
	}

	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}

	svc.Client = client
}

func (svc *KubernetesClusterService) GetClusterIps() (*[]models.NodeIp, error) {
	result := make([]models.NodeIp, 0)

	nodes, err := svc.Client.CoreV1().Nodes().List(svc.Context, metav1.ListOptions{})
	if err != nil {
		ioc.Log.Error(err.Error())
		return nil, err
	}

	for _, node := range nodes.Items {
		nodeIp := models.NodeIp{}
		for _, nodeAddress := range node.Status.Addresses {
			if nodeAddress.Type == v1.NodeHostName {
				nodeIp.Hostname = nodeAddress.Address
			}

			if nodeAddress.Type == v1.NodeExternalIP {
				nodeIp.ExternalIp = models.Ip{
					Ip:   nodeAddress.Address,
					Type: models.IPV4,
				}
			}

			if nodeAddress.Type == v1.NodeInternalIP {
				nodeIp.InternalIp = models.Ip{
					Ip:   nodeAddress.Address,
					Type: models.IPV4,
				}
			}
		}
		result = append(result, nodeIp)
	}

	return &result, nil
}

func (svc *KubernetesClusterService) GetIngressIp(serviceName string, namespace string) ([]*models.ServiceIp, error) {
	result := make([]*models.ServiceIp, 0)
	services, err := svc.Client.CoreV1().Services(namespace).List(svc.Context, metav1.ListOptions{})
	if err != nil {
		ioc.Log.Error(err.Error())
		return nil, err
	}

	for _, service := range services.Items {
		if strings.EqualFold(service.Name, serviceName) {
			if len(service.Status.LoadBalancer.Ingress) > 0 {
				serviceIp := models.ServiceIp{
					Hostname: services.Items[0].Status.LoadBalancer.Ingress[0].Hostname,
					Ip:       services.Items[0].Status.LoadBalancer.Ingress[0].IP,
					Type:     models.IPV4,
				}
				result = append(result, &serviceIp)
			}
		}
	}

	return result, nil
}

func (svc *KubernetesClusterService) GetClusterNodeNames() (*[]string, error) {
	var nodeList = make([]string, 0)

	nodes, err := svc.Client.CoreV1().Nodes().List(svc.Context, metav1.ListOptions{})
	if err != nil {
		ioc.Log.Error(err.Error())
		return nil, err
	}

	for _, node := range nodes.Items {
		nodeList = append(nodeList, node.Name)
	}

	return &nodeList, nil
}
