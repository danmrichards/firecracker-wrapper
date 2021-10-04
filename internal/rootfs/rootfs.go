package rootfs

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/danmrichards/firecracker-wrapper/internal/utils"
	"github.com/danmrichards/firecracker-wrapper/internal/vm"
)

// PrepareWorkload prepares the given VM to run a workload on startup.
func PrepareWorkload(vm *vm.VM, src, dst string) error {
	if src == "" {
		return nil
	} else if dst == "" {
		return fmt.Errorf("workload has a source (%s) but no destination", src)
	}

	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}

	// We need to mount the VM root filesystem in order to change it.
	if _, err = utils.ExecCommand("mount", "-o", "loop", vm.RootFSPath, tempDir); err != nil {
		return fmt.Errorf("mount rootfs %q: %w", vm.RootFSPath, err)
	}

	// Copy workload to the VM rootfs.
	if err = utils.CopyDir(src, filepath.Join(tempDir, dst)); err != nil {
		return fmt.Errorf("copy workload: %w", err)
	}

	if _, err = utils.ExecCommand("umount", tempDir); err != nil {
		return fmt.Errorf("mount rootfs %q: %w", vm.RootFSPath, err)
	}

	return nil
}
