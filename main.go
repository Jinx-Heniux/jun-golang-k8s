package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	// kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	// /home/zhs2si/share/Projects/kubeconfig/bmlp-ops.yaml
	kubeconfig := filepath.Join(os.Getenv("HOME"), "share", "Projects", "kubeconfig", "bmlp-ops.yaml")

	// bootstrap config
	fmt.Println()
	fmt.Println("Using kubeconfig: ", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	items, err := ListNamespaces(clientset, ctx)
	if err != nil {
		log.Fatal(err)
	}
	if items == nil {
		fmt.Println("No namespaces found!")
		return
	}

	fmt.Printf("There are %d namespaces in the cluster\n", len(items))
	for _, ns := range items {
		fmt.Printf("namespace %s\n", ns.GetName())
		labels := ns.GetLabels()
		// fmt.Printf("\tlabels: %v\n", labels)
		if value, ok := labels["field.cattle.io/projectId"]; ok {
			fmt.Printf("\tfield.cattle.io/projectId: %v\n", value)
		} else {
			fmt.Println("\tfield.cattle.io/projectId:")
		}

	}
	// for i := 0; i < len(nss.Items); i++ {
	// 	fmt.Printf("namespace %s\n", nss.Items[i].GetName())
	// }

}

func ListNamespaces(clientset *kubernetes.Clientset,
	ctx context.Context) ([]v1.Namespace, error) {

	list, err := clientset.CoreV1().Namespaces().
		List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
