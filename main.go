package main

import (
	"k8s.io/kubernetes/pkg/kubelet/cadvisor"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// cadvisor
	imageFsInfoProvider := cadvisor.NewImageFsInfoProvider("remote", "/var/run/containerd/containerd.sock")
	cAdvisorIface, err := cadvisor.New(imageFsInfoProvider, "/var/lib/kubelet", []string{"/kubepods", "/system.slice/kubelet.service"}, false)
	if err != nil {
		panic(err)
	}
	err = cAdvisorIface.Start()
	if err != nil {
		panic(err)
	}

	// block forever
	wg.Add(1)
	wg.Wait()
}
