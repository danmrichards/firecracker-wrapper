#!/bin/bash
set -ex

dpkg -i /mnt/root/linux*.deb

# Set hostname and auto login.
echo 'ubuntu-focal' > /etc/hostname
passwd -d root
mkdir /etc/systemd/system/serial-getty@ttyS0.service.d/
cat <<EOF > /etc/systemd/system/serial-getty@ttyS0.service.d/autologin.conf
[Service]
ExecStart=
ExecStart=-/sbin/agetty --autologin root -o '-p -- \\u' --keep-baud 115200,38400,9600 %I $TERM
EOF

# Enable the manager service
chmod +x /opt/manager-linux-amd64
systemctl enable manager.service

# Enable the init-network service
chmod +x /opt/init-network.sh
systemctl enable init-network.service
echo "nameserver 8.8.8.8" >> /etc/resolv.conf

