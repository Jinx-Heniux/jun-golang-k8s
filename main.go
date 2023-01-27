package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	api := clientset.CoreV1()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println(api)
	// fmt.Println(ctx)

	nss, err := api.Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(nss)

	if len(nss.Items) == 0 {
		fmt.Println("No namespaces found!")
		return
	}

	fmt.Printf("There are %d namespaces in the cluster\n", len(nss.Items))
	for _, ns := range nss.Items {
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
