package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	internalapi "k8s.io/cri-api/pkg/apis"
	"k8s.io/kubernetes/pkg/kubelet/cri/remote"
	"net"
	"sync"
	"time"

	"github.com/hanlins/xray"
	"github.com/hanlins/xray/dot"
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

	s := xray.NewScanner(nil)
	nodeCh := s.Scan(rs)
	g, _ := dot.Draw(dot.NewGraphInfo(s), nodeCh, nil)
	fmt.Printf("dot object:\n%s\n", g)
}

func dialer(ctx context.Context, addr string) (net.Conn, error) {
	return (&net.Dialer{}).DialContext(ctx, "unix", addr)
}

func dialEp(endpoint string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var maxMsgSize = 1024 * 1024 * 16
	_, err := grpc.DialContext(ctx, endpoint, grpc.WithInsecure(), grpc.WithContextDialer(dialer), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))
	return err
}

func main() {
	var containerID string
	var endpoint string
	var interval int64
	var verbose bool
	var rtStatus bool
	var parallel int

	flag.StringVar(&containerID, "id", "fake", "container id")
	flag.StringVar(&endpoint, "ep", "/var/run/containerd/containerd.sock", "container runtime endpoint")
	flag.Int64Var(&interval, "int", 5000, "loop interval in ms")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.BoolVar(&rtStatus, "rt", false, "poll container runtime status instead of specific container's status")
	flag.IntVar(&parallel, "p", 1, "number of parallel goroutines polling status")

	flag.Parse()

	var wg sync.WaitGroup

	rs, err := remote.NewRemoteRuntimeService(endpoint, time.Minute)
	if err != nil {
		fmt.Printf("NewRemoteRuntimeService err: %+v\n", err)
		return
	}

	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func() {
			// loop for container
			for {
				err := dialEp(endpoint)
				if err != nil {
					fmt.Printf("diap endpoint '%s' error: %v\n", endpoint, err)
				} else {
					fmt.Printf("diap endpoint '%s' success\n", endpoint)
				}
				if rtStatus {
					checkRuntimeStatus(rs, verbose)
				} else {
					checkContainerStatus(rs, containerID, verbose)
				}
				time.Sleep(time.Millisecond * time.Duration(interval))
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
