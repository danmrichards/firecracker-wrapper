#!/bin/bash
set -ex

rm -rf /output/*

cp /root/linux-source-$KERNEL_SOURCE_VERSION/vmlinux /output/vmlinux
cp /root/linux-source-$KERNEL_SOURCE_VERSION/.config /output/config

truncate -s 1G /output/image.ext4
mkfs.ext4 /output/image.ext4

mount /output/image.ext4 /rootfs
debootstrap --include curl,openssh-server,unzip,rsync,apt,nano focal /rootfs http://archive.ubuntu.com/ubuntu/

cp /opt/manager-linux-amd64 /rootfs/opt/manager-linux-amd64
cp /script/manager.service /rootfs/etc/systemd/system/manager.service

cp /script/init-network.sh /rootfs/opt/init-network.sh
cp /script/init-network.service /rootfs/etc/systemd/system/init-network.service

mount --bind / /rootfs/mnt
chroot /rootfs /bin/bash /mnt/script/provision.sh
