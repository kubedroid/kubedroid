#!/bin/bash
setfacl -m u:qemu:rw /dev/dri/renderD128

getfacl /dev/dri/renderD128
cat /etc/libvirt/qemu.conf
