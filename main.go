package main

import (
	"flag"
	"fmt"
	"k8s.io/kubernetes/pkg/kubelet/cri/remote"
	"time"
)

func main() {
	var containerID string
	var endpoint string
	var interval int64
	var verbose bool

	flag.StringVar(&containerID, "id", "fake", "container id")
	flag.StringVar(&endpoint, "ep", "/var/run/containerd/containerd.sock", "container runtime endpoint")
	flag.Int64Var(&interval, "int", 5000, "loop interval in ms")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	rs, err := remote.NewRemoteRuntimeService(endpoint, time.Minute)
	if err != nil {
		fmt.Printf("NewRemoteRuntimeService err: %+v\n", err)
		return
	}

	// loop for container
	for {
		status, err := rs.ContainerStatus(containerID)
		if err != nil {
			fmt.Printf("ContainerStatus err: %+v\n", err)
		} else {
			if verbose {
				fmt.Printf("ContainerStatus %s status: %+v\n", containerID, status)
			} else {
				fmt.Println(containerID + " status ok")
			}
		}
		time.Sleep(time.Millisecond * time.Duration(interval))
	}
}
