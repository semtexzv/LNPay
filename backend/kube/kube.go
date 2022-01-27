package kube

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var clientset *kubernetes.Clientset

func init() {
	config, err := rest.InClusterConfig()
	if err != nil {
		return
	}
	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	println("External IP", GetSvcEndpoint("lnd-external"))
}

func GetSvcEndpoint(name string) string {
	if clientset == nil {
		return ""
	}
	svc, err := clientset.CoreV1().Services("default").
		Get(context.Background(),name, metav1.GetOptions{})

	if err != nil {
		panic(err)
		return ""
	}
	for _, ing := range svc.Status.LoadBalancer.Ingress {
		println("IP", ing.IP, ing.Hostname)
		return ing.IP
	}
	return ""
}
