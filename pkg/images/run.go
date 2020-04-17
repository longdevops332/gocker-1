package images

import (
	"fmt"
	"os"

	"github.com/containerd/cgroups"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// Run defines struct for running container
type Run struct {
}

// NewRun provides starting of the new container
func NewRun(name string) (*Run, error) {
	baseDir := os.Getenv("GOCKER_BASE_DIR")
	if baseDir == "" {
		baseDir = "gocker-images"
	}
	return &Run{}, nil
}

// Do provides starting of container
func (r *Run) Do() error {
	shares := uint64(100)
	control, err := cgroups.New(cgroups.V1, cgroups.StaticPath("/test"), &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Shares: &shares,
		},
	})
	if err != nil {
		return fmt.Errorf("unable to create cgroup: %v", err)
	}
	defer control.Delete()
	return nil
}
