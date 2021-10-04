# Firecracker Wrapper
A quick and dirty wrapper around [Firecracker][1] to demonstrate the way in which
microvms can be created on the fly.

## Requirements
* Go 1.13+
* Docker
* [Firecracker][2]

## Building From Source

Clone this repo and build the binaries:

```bash
$ make build
```

## Dependencies
### Go
Update the Go dependencies like so:

```bash
$ make deps
```

### Submodules
This application relies on the [Ubuntu Firecracker][3] project which provides a quick
and easy way to create a kernel and root filesystem image for Firecracker. This
dependency is managed via a Git submodule, you can create/update it like so:

```bash
$ make ubuntu-firecracker
```

## Usage

First ensure that you have built the kernel and rootfs images like so:

```bash
$ make images
```

Then you can use the application:
```bash

```

[1]: https://firecracker-microvm.github.io/
[2]: https://github.com/firecracker-microvm/firecracker#getting-started
[3]: https://github.com/bkleiner/ubuntu-firecracker
