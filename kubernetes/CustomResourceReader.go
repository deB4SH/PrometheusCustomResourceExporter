package kubernetes

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type CustomResource struct {
	Api        string
	ApiVersion string
	Namespace  string

	Resource string
	Name     string
}

type CustomResourceData struct {
	Name string
	Data string
}

type KubernetesConnection struct {
	config *rest.Config
}

var isCluster = flag.Bool("is.cluster", false, "Boolean toggle to define if application is deployed inside a cluster or not.")

func NewCustomResource(api string, apiversion string, namespace string, resource string, name string) *CustomResource {
	e := new(CustomResource)
	e.Api = api
	e.ApiVersion = apiversion
	e.Namespace = namespace
	e.Resource = resource
	e.Name = name
	return e
}

func NewCustomResourceData(name string, data []byte) *CustomResourceData {
	e := new(CustomResourceData)
	e.Name = name
	e.Data = string(data)
	return e
}

func newKubernetesConnection(config *rest.Config) *KubernetesConnection {
	e := new(KubernetesConnection)
	e.config = config
	return e
}

func BuildKubernetesConnection() *KubernetesConnection {
	var k8Config *KubernetesConnection

	if *isCluster {
		fmt.Print("Starting to use incluster config")
		config, err := rest.InClusterConfig()
		if err != nil {
			fmt.Printf("error getting Kubernetes config: %v\n", err)
			os.Exit(1)
		}
		k8Config = newKubernetesConnection(config)
	} else {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("error getting user home dir: %v\n", err)
			os.Exit(1)
		}
		kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			fmt.Printf("error getting Kubernetes config: %v\n", err)
			os.Exit(1)
		}
		k8Config = newKubernetesConnection(config)
	}

	return k8Config
}

func ParseCR(c *CustomResource, k8Config KubernetesConnection) *CustomResourceData {
	fmt.Print("Creating clientset")
	clientset, err := kubernetes.NewForConfig(k8Config.config)
	if err != nil {
		fmt.Printf("Could not aquire clientset from config: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("Query requested cr")
	data, err := clientset.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/%s/%s", c.Api, c.ApiVersion)).Namespace(c.Namespace).Resource(c.Resource).Name(c.Name).DoRaw(context.TODO())
	if err != nil {
		fmt.Printf("Could not find custom resource:\n %v\n", err)
		os.Exit(1)
	}

	return NewCustomResourceData(fmt.Sprintf("%s-%s", c.Resource, c.Name), data)
}
