package images

import (
	"fmt"
	"os"

	"github.com/containerd/cgroups"
	"github.com/google/uuid"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
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
	id := genID()
	name := fmt.Sprintf("c_%s", id)
	logrus.Infof("Prepare to start container with id %s", name)
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

// genID provides generation of unique id
func genID() string {
	return uuid.New().String()
}
