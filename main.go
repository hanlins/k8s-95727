package main

import (
	"flag"
	"fmt"
	internalapi "k8s.io/cri-api/pkg/apis"
	"k8s.io/kubernetes/pkg/kubelet/cri/remote"
	"time"
)

func checkContainerStatus(rs internalapi.RuntimeService, containerID string, verbose bool) {
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
}

func checkRuntimeStatus(rs internalapi.RuntimeService, verbose bool) {
	status, err := rs.Status()
	if err != nil {
		fmt.Printf("Container Runtime status err: %+v\n", err)
	} else {
		if verbose {
			fmt.Printf("Container Runtime status: %+v\n", status)
		} else {
			fmt.Println("runtime status ok")
		}
	}
}

func main() {
	var containerID string
	var endpoint string
	var interval int64
	var verbose bool
	var rtStatus bool

	flag.StringVar(&containerID, "id", "fake", "container id")
	flag.StringVar(&endpoint, "ep", "/var/run/containerd/containerd.sock", "container runtime endpoint")
	flag.Int64Var(&interval, "int", 5000, "loop interval in ms")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.BoolVar(&rtStatus, "rt", false, "poll container runtime status instead of specific container's status")

	flag.Parse()

	rs, err := remote.NewRemoteRuntimeService(endpoint, time.Minute)
	if err != nil {
		fmt.Printf("NewRemoteRuntimeService err: %+v\n", err)
		return
	}

	// loop for container
	for {
		if rtStatus {
			checkRuntimeStatus(rs, verbose)
		} else {
			checkContainerStatus(rs, containerID, verbose)
		}
		time.Sleep(time.Millisecond * time.Duration(interval))
	}
}
