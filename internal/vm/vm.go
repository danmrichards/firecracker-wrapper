package vm

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/danmrichards/firecracker-wrapper/internal/config"
	"github.com/danmrichards/firecracker-wrapper/internal/utils"
	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
)

const dataDir = "data"

// VM represents a Firecracker VM.
type VM struct {
	dataDir    string
	RootFSPath string
	Config     firecracker.Config
}

// New returns a configured Firecracker VM.
func New(id string, cfg *config.WrapperConfig) (*VM, error) {
	// Each VM has a unique data directory.
	vmDataDir := filepath.Join(dataDir, id)
	if err := os.MkdirAll(vmDataDir, 0755); err != nil {
		return nil, fmt.Errorf("create vm data dir: %w", err)
	}

	vmKernelImg := filepath.Join(vmDataDir, filepath.Base(cfg.KernelImg))
	vmRootFS := filepath.Join(vmDataDir, filepath.Base(cfg.RootFS))

	if err := initVM(cfg.KernelImg, vmKernelImg, cfg.RootFS, vmRootFS); err != nil {
		return nil, fmt.Errorf("init vm: %w", err)
	}

	sp, err := socketPath(id)
	if err != nil {
		return nil, fmt.Errorf("socket path: %w", err)
	}

	return &VM{
		dataDir:    vmDataDir,
		RootFSPath: vmRootFS,
		Config: firecracker.Config{
			SocketPath:      sp,
			VMID:            id,
			KernelImagePath: vmKernelImg,
			KernelArgs:      cfg.KernelOpts,
			Drives: []models.Drive{
				{
					DriveID:      firecracker.String("rootfs"),
					PathOnHost:   firecracker.String(vmRootFS),
					IsRootDevice: firecracker.Bool(true),
					IsReadOnly:   firecracker.Bool(false),
				},
			},
			MachineCfg: models.MachineConfiguration{
				HtEnabled:  firecracker.Bool(false),
				MemSizeMib: firecracker.Int64(cfg.Memory),
				VcpuCount:  firecracker.Int64(cfg.NumCPUs),
			},
			DisableValidation: false,
		},
	}, nil
}

// Clean cleans up the VM.
func (vm *VM) Clean() error {
	return os.RemoveAll(vm.dataDir)
}

func initVM(srcKernelImg, dstKernelImg, srcRootFS, dstRootFS string) (err error) {
	if err = utils.CopyFile(srcKernelImg, dstKernelImg); err != nil {
		return fmt.Errorf("kernel image: %w", err)
	}
	if err = utils.CopyFile(srcRootFS, dstRootFS); err != nil {
		return fmt.Errorf("rootfs: %w", err)
	}

	return nil
}

// socketPath provides a randomized socket path by building a unique filename
// and searching for the existence of directories {$HOME, os.TempDir()} and
// returning the path with the first directory joined with the unique filename.
// If we can't find a good path panics.
func socketPath(id string) (string, error) {
	filename := strings.Join([]string{".firecracker.sock", id}, "-")
	var dir string
	if d := os.Getenv("HOME"); utils.DirectoryExists(d) {
		dir = d
	} else if utils.DirectoryExists(os.TempDir()) {
		dir = os.TempDir()
	} else {
		return "", errors.New("could not find location for socket")
	}

	return filepath.Join(dir, filename), nil
}
