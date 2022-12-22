package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	kubeconfig := filepath.Join(userHomeDir, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	go deletePod(clientset)

	// wait forever
	select {}
}

func deletePod(clientset *kubernetes.Clientset) {
	fmt.Println("Deleting pod")
	for {
		// get all pods
		pods, err := clientset.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		// select one random pod
		pod := pods.Items[rand.Intn(len(pods.Items))]
		// delete pod
		err = clientset.CoreV1().Pods("default").Delete(context.Background(), pod.Name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Deleted pod:", pod.Name)
		time.Sleep(1 * time.Second)
	}
}

func deleteNode(clientset *kubernetes.Clientset) {
	fmt.Println("Deleting node")
	for {
		// get all nodes
		nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		// select one random node
		node := nodes.Items[rand.Intn(len(nodes.Items))]
		// delete node
		err = clientset.CoreV1().Nodes().Delete(context.Background(), node.Name, metav1.DeleteOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Deleted node:", node.Name)
		time.Sleep(300 * time.Second)
	}
}
