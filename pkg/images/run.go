package images

import (
	"fmt"
	"os"

	"github.com/containerd/cgroups"
	"github.com/google/uuid"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

// Run defines struct for running container
type Run struct {
	imageName, deviceName string
}

// NewRun provides starting of the new container
func NewRun(name, deviceName string) (*Run, error) {
	baseDir := os.Getenv("GOCKER_BASE_DIR")
	if baseDir == "" {
		baseDir = "gocker-images"
	}
	return &Run{
		imageName:  name,
		deviceName: deviceName,
	}, nil
}

// Do provides starting of container
func (r *Run) Do() error {
	id := genID()
	fmt.Println(r.deviceName)
	name := fmt.Sprintf("c_%s", id)
	logrus.Infof("Prepare to start container with id %s", name)
	path := fmt.Sprintf("gocker/%s", r.imageName)
	shares := uint64(100)
	control, err := cgroups.New(cgroups.V1, cgroups.StaticPath(path), &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Shares: &shares,
		},
	})
	if err := control.Add(cgroups.Process{
		Pid: os.Getpid(),
	}); err != nil {
		return fmt.Errorf("unable to add process to cgroup: %v", err)
	}
	defer control.Delete()
	if err != nil {
		return fmt.Errorf("unable to create cgroup: %v", err)
	}

	if err := createNetwork(r.imageName, r.deviceName); err != nil {
		return fmt.Errorf("unable to create network: %v", err)
	}
	return nil
}

// genID provides generation of unique id
func genID() string {
	return uuid.New().String()
}

// createNetwork provides creating of the new network on container
func createNetwork(name, networkName string) error {
	la := netlink.NewLinkAttrs()
	la.Name = name
	mybridge := &netlink.Bridge{LinkAttrs: la}
	err := netlink.LinkAdd(mybridge)
	if err != nil {
		return fmt.Errorf("unable to add link: %s %v", name, err)
	}
	eth1, err := netlink.LinkByName(networkName)
	if err != nil {
		return fmt.Errorf("failed to define netork: %v", err)
	}
	netlink.LinkSetMaster(eth1, mybridge)
	return nil
}
