package firecracker

import (
	"context"
	"fmt"
	"os"

	"github.com/danmrichards/firecracker-wrapper/internal/vm"
	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/sirupsen/logrus"
)

const executableMask = 0111

// Wrapper wraps Firecracker.
type Wrapper struct {
	bin string
}

// NewWrapper returns a new wrapper.
func NewWrapper(bin string) *Wrapper {
	return &Wrapper{
		bin: bin,
	}
}

// RunVM runs the given VM.
func (w *Wrapper) RunVM(ctx context.Context, logger *logrus.Logger, vm *vm.VM) error {
	vmmCtx, vmmCancel := context.WithCancel(ctx)
	defer vmmCancel()

	machineOpts := []firecracker.Opt{
		firecracker.WithLogger(logrus.NewEntry(logger)),
	}

	if finfo, err := os.Stat(w.bin); os.IsNotExist(err) {
		return fmt.Errorf("firecracker binary %q does not exist: %v", w.bin, err)
	} else if err != nil {
		return fmt.Errorf("failed to stat binary, %q: %v", w.bin, err)
	} else if finfo.IsDir() {
		return fmt.Errorf("binary, %q, is a directory", w.bin)
	} else if finfo.Mode()&executableMask == 0 {
		return fmt.Errorf("binary, %q, is not executable. Check permissions of binary", w.bin)
	}

	cmd := firecracker.VMCommandBuilder{}.
		WithBin(w.bin).
		WithSocketPath(vm.Config.SocketPath).
		WithStdin(os.Stdin).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		Build(ctx)

	machineOpts = append(machineOpts, firecracker.WithProcessRunner(cmd))

	m, err := firecracker.NewMachine(vmmCtx, vm.Config, machineOpts...)
	if err != nil {
		return fmt.Errorf("create machine: %w", err)
	}

	if err = m.Start(vmmCtx); err != nil {
		return fmt.Errorf("start machine: %w", err)
	}
	defer m.StopVMM()

	// Wait for the VM to exit.
	if err = m.Wait(vmmCtx); err != nil {
		return fmt.Errorf("vm wait: %w", err)
	}

	return nil
}
