package images

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/containerd/cgroups"
	"github.com/google/uuid"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/saromanov/gocker/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

// request defines helpful struct for running container
type request struct {
	workingDir string
	cmd        string
	path       string
	shares     uint64
}

// Run defines struct for running container
type Run struct {
	imageName, deviceName, baseDir string
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
		baseDir:    baseDir,
	}, nil
}

// Do provides starting of container
func (r *Run) Do() error {
	id := genID()
	fmt.Println(r.deviceName)
	name := fmt.Sprintf("c_%s", id)
	logrus.Infof("Prepare to start container with id %s", name)
	path := r.prepareImagePath(r.imageName)
	manifest, err := loadManifest(fmt.Sprintf("%s.json", path))
	if err != nil {
		return err
	}
	state, err := manifestHistoryToJSON(manifest.History[0].V1Compatibility)
	if err != nil {
		return err
	}
	shares := uint64(100)
	return r.run(request{
		path:       path,
		shares:     shares,
		workingDir: state.Config.WorkingDir,
		cmd:        strings.Join(state.Config.Cmd, " "),
	})
}

func (r *Run) run(req request) error {
	control, err := cgroups.New(cgroups.V1, cgroups.StaticPath(req.path), &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Shares: &req.shares,
		},
	})

	if err != nil {
		return fmt.Errorf("unable to create cgroup: %v", err)
	}
	defer control.Delete()

	if err := control.Add(cgroups.Process{
		Pid:  os.Getpid(),
		Path: req.path,
	}); err != nil {
		return fmt.Errorf("unable to add process to cgroup: %v", err)
	}
	chExit, err := chroot(req.path)
	if err != nil {
		return fmt.Errorf("unable to make chroot of dir %s: %v", req.path, err)
	}
	if req.workingDir != "" {
		if err := os.Chdir(req.workingDir); err != nil {
			return fmt.Errorf("unable to chdir: %s %v", req.workingDir, err)
		}
	}
	defer func() {
		if err := chExit(); err != nil {
			panic(err)
		}
	}()

	/*if err := createNetwork(r.imageName, r.deviceName); err != nil {
		return fmt.Errorf("unable to create network: %v", err)
	}*/

	return execCmd(req.workingDir, req.cmd)
}

func (r *Run) prepareImagePath(name string) string {
	if !strings.Contains(name, "/") {
		return fmt.Sprintf("%s/library_%s", r.baseDir, name)
	}
	splitting := strings.Split(name, "/")
	if len(splitting) <= 1 {
		return fmt.Sprintf("library")
	}
	return fmt.Sprintf("%s/%s_%s", r.baseDir, splitting[0], splitting[1])
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
		return fmt.Errorf("failed to define network: %v", err)
	}
	netlink.LinkSetMaster(eth1, mybridge)
	return nil
}

// chroot provides setting chroot of teh dir
func chroot(path string) (func() error, error) {
	root, err := os.Open("/")
	if err != nil {
		return nil, err
	}

	if err := syscall.Chroot(path); err != nil {
		root.Close()
		return nil, err
	}

	return func() error {
		defer root.Close()
		if err := root.Chdir(); err != nil {
			return err
		}
		return syscall.Chroot(".")
	}, nil
}

// loading of the manifest for running image
func loadManifest(path string) (*models.Manifest, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to load manifest file: %s %v", path, err)
	}
	m := &models.Manifest{}
	if err := json.Unmarshal([]byte(file), &m); err != nil {
		return nil, fmt.Errorf("unable to unmarshal manifest file: %v", err)
	}
	return m, nil
}

func manifestHistoryToJSON(s string) (*models.V1Compatibility, error) {
	conf := &models.V1Compatibility{}
	if err := json.Unmarshal([]byte(s), &conf); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config: %v", err)
	}
	return conf, nil
}

func execCmd(path, cmdPath string) error {
	fmt.Println(cmdPath)
	cmdGoVer := &exec.Cmd{
		Path:   cmdPath,
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	if err := cmdGoVer.Run(); err != nil {
		return fmt.Errorf("unable to execute command: %v", err)
	}
	return nil
}
