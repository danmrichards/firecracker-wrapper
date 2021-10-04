package main

import (
	"context"
	"flag"
	"log"

	"github.com/danmrichards/firecracker-wrapper/internal/config"
	"github.com/danmrichards/firecracker-wrapper/internal/firecracker"
	"github.com/danmrichards/firecracker-wrapper/internal/rootfs"
	"github.com/danmrichards/firecracker-wrapper/internal/vm"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	cfg := config.New()
	cfg.BindFlags(flag.CommandLine)
	flag.Parse()

	if err := cfg.Validate(); err != nil {
		logger.Fatal(err)
	}

	if cfg.WorkloadSrc == "" {
		logger.Infof("No workload specified")
	}

	vmID := uuid.NewString()
	fvm, err := vm.New(vmID, cfg)
	if err != nil {
		logger.Fatal(err)
	}

	if err = rootfs.PrepareWorkload(fvm, cfg.WorkloadSrc, cfg.WorkloadDst); err != nil {
		logger.Fatal(err)
	}

	logger.Infof("VM %q launching", vmID)
	fw := firecracker.NewWrapper(cfg.FirecrackerBin)

	if err = fw.RunVM(context.Background(), logger, fvm); err != nil {
		log.Fatal(err)
	}

	logger.Infof("VM %q exited", vmID)
	if err = fvm.Clean(); err != nil {
		logger.Fatal(err)
	}
}
