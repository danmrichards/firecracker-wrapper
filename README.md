# Firecracker Wrapper
A quick and dirty wrapper around [Firecracker][1] to demonstrate the way in which
microvms can be created on the fly. I am aware that [firectl][2] exists, this
is just a learning exercise.

This repo consists of two applications, the wrapper and the manager. The manager
provides a lightweight REST API and acts as a means of starting a process inside
the VM after it has booted. The wrapper is what it says on the tin, it wraps the
firecracker binary to start the VM.

## Requirements
* Go 1.13+
* Docker
* [Firecracker][3]

## Building From Source

Clone this repo and build the binaries:

```bash
$ make build
```

Ensure you also build the kernel and rootfs images:

```bash
$ docker build -t firecracker-manager -f manager/Dockerfile . && \
  sudo docker run --privileged -it --rm -v $(pwd)/manager/output:/output firecracker-manager
```

## Dependencies
### Go
Update the Go dependencies like so:

```bash
$ make deps
```

## Usage

First you will need to configure a tap network device on your host:

```bash
$ make tap
```

```bash
Usage of ./bin/wrapper-linux-amd64:
  -firecracker-binary string
    	Path to firecracker binary (default "firecracker")
  -kernel-image string
    	Path to the kernel image
  -kernel-opts string
    	Kernel commandline (default "ro console=ttyS0 noapic reboot=k panic=1 pci=off nomodules")
  -memory int
    	VM memory, in MiB (default 512)
  -numcpus int
    	Number of CPUs for the VM (default 1)
  -rootfs string
    	Path to root disk image
  -workload-dst string
    	Path to the directory where the workload should be created
  -workload-exe string
    	Path to the workload executable
  -workload-src string
    	Path to the directory containing the workload
```

When you're done, remove the tap network device:

```bash
$ make cleantap
```

[1]: https://firecracker-microvm.github.io/
[2]: https://github.com/firecracker-microvm/firectl
[3]: https://github.com/firecracker-microvm/firecracker#getting-started
