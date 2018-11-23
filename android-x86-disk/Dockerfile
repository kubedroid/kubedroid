# This Dockerfile will create a RegistryDisk which contains
# android-x86. You can use it to provision an Android-x86
# VM in KubeVirt.
#
# See https://kubevirt.io/user-guide/#/workloads/virtual-machines/disks-and-volumes?id=registrydisk
# for more information about RegistryDisks
#
# The contents is sourced from the android-x86-base image,
# which contain the kernel, uncompressed ramdisk, uncompressed
# initrd, and the base system.
#
# The build process will:
# - Generate the compressed ramdisk and initrd images
# - Create a target file system, based of the
#   file system in the system container.
# - Add the kernel, ramdisk and initrd images to the
#   target file system
# - Create a raw virtual disk, with a single partition
# - Copy the target file system to that partition
# - Install GRUB
# - Create a RegistryDisk container which contains a
#   single qcow2 image.
#
# To create the initrd and ramdisk images, mkbootfs
# is used.
#
# To install GRUB in an offline disk without root
# privileges (docker build doesn't have any privileges),
# a patched version of GRUB is used.

# -------------------------------------------------
# Fetch the base images. They are based on the .iso
# files which Android-x86 releases.
# -------------------------------------------------

FROM quay.io/quamotion/android-x86-base:8.1-rc2 AS base

ENV image_name=android-x86

WORKDIR /android
# Apply the patches
# COPY *.patch .
# RUN apt-get update \
# && apt-get install -y patch \
# && echo system/build.prop \
# && ls -l *.patch \
# && patch -p1 < android-x86.patch \
# && echo system/build.prop \
# && rm *.patch

# Update the ramdisk and initrd images
RUN mkbootfs ./ramdisk | gzip > ramdisk.img \
&& rm -rf ./ramdisk \
&& mkbootfs ./initrd | gzip > initrd.img \
&& rm -rf ./initrd

# Inject the grub configuration
COPY grub.cfg grub/grub.cfg

# Build the GRUB bootloader
COPY grub-early.cfg /tmp/grub-early.cfg
COPY android-x86.sfdisk /tmp/android-x86.sfdisk

# Prepare the GRUB files in /android/grub2
ENV GRUB2_MODULES="biosdisk boot chain configfile ext2 linux ls part_msdos reboot serial vga"

RUN mkdir -p /tmp/grub2/ \
&& which grub-mkimage \
&& which grub-bios-setup \
&& /usr/local/bin/grub-mkimage \
        -p /boot/grub \
        -d /usr/local/lib/grub/i386-pc \
        -o /tmp/grub2/core.img \
        -O i386-pc \
        -c /tmp/grub-early.cfg \
        ${GRUB2_MODULES} \
&& cp /usr/local/lib/grub/i386-pc/*.img /tmp/grub2/ \
&& echo "(hd0) /tmp/${image_name}.img" > /tmp/device.map

# - Create a 4GB raw image
# - Partition it with one 4GB - 1MB partition
# - Format the partion with ext2fs
# - Copy the contents of the rootfs folder to that partition
# - Install GRUB2
# - Convert it to qcow2
ENV image_size=4096
ENV block_size=1024

RUN dd if=/dev/zero of=/tmp/$image_name.img bs=$block_size count=$(( $image_size * 1024)) \
&& printf " \
type=83, size=$(( ($image_size - 1) * 1024 * 2)) \
" | sfdisk /tmp/$image_name.img \
&& fdisk -l /tmp/$image_name.img \
&& mke2fs \
	-t ext4 \
	-F \
	-d /android/ \
        # 1024 bytes per block
	-b $block_size \
        # label
	-L android \
	/tmp/$image_name.img \
	-E offset=$((2048 * 512)) \
        # blocks count, 1024 * block size (4096) = 4 GB,
        # but there's a 1 MB boot sector
	$(( (image_size - 1) * 1024)) \
&& /usr/local/sbin/grub-bios-setup \
        --device-map=/tmp/device.map \
        -d /tmp/grub2/ \
        /tmp/$image_name.img \
	-r "hd0,msdos1" \
&& qemu-img convert -f raw -O qcow2 /tmp/$image_name.img /tmp/$image_name.qcow2 \
&& rm /tmp/$image_name.img

FROM kubevirt/registry-disk-v1alpha

COPY --from=base /tmp/*.qcow2 /disk/
