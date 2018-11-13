# This Dockerfile will create a RegistryDisk which contains
# android-x86. You can use it to provision an Android-x86
# VM in KubeVirt.
#
# See https://kubevirt.io/user-guide/#/workloads/virtual-machines/disks-and-volumes?id=registrydisk
# for more information about RegistryDisks
#
# The contents is sourced from four containers, which contain
# the kernel, uncompressed ramdisk, uncompressed initrd,
# and the base system.
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
# is used. A container which contains mkbootfs is
# available.
#
# To install GRUB in an offline disk without root
# privileges (docker build doesn't have any privileges),
# a patched version of GRUB is used.

# -------------------------------------------------
# Fetch the base images. They are based on the .iso
# files which Android-x86 releases.
# -------------------------------------------------

FROM kubevirt/registry-disk-v1alpha
FROM quamotion/android-x86-kernel:7.1-r2 AS kernel

FROM quamotion/android-x86-ramdisk:7.1-r2 AS ramdisk

FROM quamotion/android-x86-initrd:7.1-r2 AS initrd

FROM quamotion/android-x86-system:7.1-r2 AS system

# -------------------------------------------------
# Make the initrd and ramdisk images.
# -------------------------------------------------
FROM quamotion/android-x86-mkbootfs:latest AS mkbootfs

WORKDIR /android/initrd
COPY --from=initrd / .

WORKDIR /android/ramdisk
COPY --from=ramdisk / .

# Update the ramdisk and initrd images
WORKDIR /android
RUN mkdir -p rootfs \
&& mkbootfs ./ramdisk | gzip > rootfs/ramdisk.img \
&& mkbootfs ./initrd | gzip > rootfs/initrd.img

# -------------------------------------------------
# Build the target rootfs, and generate the qcow2
# image
# -------------------------------------------------
FROM quamotion/android-x86-build AS build

# Copy the kernel and system images "as is",
# extract the compressed initrd and ramdisk images

ENV image_name=android-x86

WORKDIR /android/rootfs
COPY --from=kernel /kernel .
COPY --from=mkbootfs /android/rootfs/ramdisk.img ramdisk.img
COPY --from=mkbootfs /android/rootfs/initrd.img initrd.img

WORKDIR /android/rootfs/system
COPY --from=system / .

WORKDIR /android
COPY grub.cfg /android/rootfs/grub/grub.cfg
COPY android-x86.sfdisk .

# Prepare the GRUB files in /android/grub2
ENV GRUB2_MODULES="biosdisk boot chain configfile ext2 linux ls part_msdos reboot serial vga"
COPY grub-early.cfg .

RUN mkdir -p grub2/ \
&& which grub-mkimage \
&& which grub-bios-setup \
&& /usr/local/bin/grub-mkimage \
        -p /boot/grub \
        -d /usr/local/lib/grub/i386-pc \
        -o grub2/core.img \
        -O i386-pc \
        -c ./grub-early.cfg \
        ${GRUB2_MODULES} \
&& cp /usr/local/lib/grub/i386-pc/*.img grub2/ \
&& echo "(hd0) /android/${image_name}.img" > device.map

# - Create a 4GB raw image
# - Partition it with one 4GB - 1MB partition
# - Format the partion with ext2fs
# - Copy the contents of the rootfs folder to that partition
# - Install GRUB2
# - Convert it to qcow2
ENV image_size=4096
ENV block_size=1024

RUN dd if=/dev/zero of=/android/$image_name.img bs=$block_size count=$(( $image_size * 1024)) \
&& printf " \
type=83, size=$(( ($image_size - 1) * 1024 * 2)) \
" | sfdisk /android/$image_name.img \
&& fdisk -l /android/$image_name.img \
&& mke2fs \
	-F \
	-d /android/rootfs \
        # 1024 bytes per block
	-b $block_size \
        # label
	-L android \
	$image_name.img \
	-E offset=$((2048 * 512)) \
        # blocks count, 1024 * block size (4096) = 4 GB,
        # but there's a 1 MB boot sector
	$(( (image_size - 1) * 1024)) \
&& /usr/local/sbin/grub-bios-setup \
        --device-map=device.map \
        -d grub2/ \
        /android/$image_name.img \
	-r "hd0,msdos1" \
&& qemu-img convert -f raw -O qcow2 /android/$image_name.img /android/$image_name.qcow2 \
&& rm /android/$image_name.img

FROM kubevirt/registry-disk-v1alpha

COPY --from=build /android/*.qcow2 /disk/
