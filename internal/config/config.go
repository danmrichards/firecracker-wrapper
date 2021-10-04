package config

import (
	"errors"
	"flag"
)

type WrapperConfig struct {
	FirecrackerBin string
	KernelImg      string
	KernelOpts     string
	RootFS         string
	RootSize       int
	NumCPUs        int64
	Memory         int64
	WorkloadSrc    string
	WorkloadDst    string
	WorkloadExe    string
}

// New returns an empty WrapperConfig with zero values
func New() *WrapperConfig {
	return &WrapperConfig{}
}

// BindFlags parses the given flag.FlagSet and sets the options accordingly.
func (w *WrapperConfig) BindFlags(fs *flag.FlagSet) {
	fs.StringVar(&w.FirecrackerBin, "firecracker-binary", "firecracker", "Path to firecracker binary")
	fs.StringVar(&w.KernelImg, "kernel-image", "", "Path to the kernel image")
	fs.StringVar(&w.KernelOpts, "kernel-opts", "ro console=ttyS0 noapic reboot=k panic=1 pci=off nomodules", "Kernel commandline")
	fs.StringVar(&w.RootFS, "rootfs", "", "Path to root disk image")
	fs.IntVar(&w.RootSize, "rootfs-size", 5, "Size, in GiB, for the root filesystem")
	fs.Int64Var(&w.NumCPUs, "numcpus", 1, "Number of CPUs for the VM")
	fs.Int64Var(&w.Memory, "memory", 512, "VM memory, in MiB")
	fs.StringVar(&w.WorkloadSrc, "workload-src", "", "Path to the directory containing the workload")
	fs.StringVar(&w.WorkloadDst, "workload-dst", "", "Path to the directory where the workload should be created")
	fs.StringVar(&w.WorkloadExe, "workload-exe", "", "Path to the workload executable")
}

// Validate returns an error if the configuration is not valid.
func (w *WrapperConfig) Validate() error {
	if w.KernelImg == "" {
		return errors.New("kernel-image is required")
	}
	if w.RootFS == "" {
		return errors.New("rootfs is required")
	}

	return nil
}
