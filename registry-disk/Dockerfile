FROM quamotion/android-x86-kernel:7.1-r2 AS kernel

FROM quamotion/android-x86-ramdisk:7.1-r2 AS ramdisk

FROM quamotion/android-x86-initrd:7.1-r2 AS initrd

FROM quamotion/android-x86-system:7.1-r2 AS system

FROM quamotion/android-x86-build:latest AS compressed

WORKDIR /android/initrd
COPY --from=initrd / .

WORKDIR /android/ramdisk
COPY --from=ramdisk / .

# Update the ramdisk and initrd images
WORKDIR /android
RUN mkdir -p rootfs \
&& mkbootfs ./ramdisk | gzip > rootfs/ramdisk.img \
&& mkbootfs ./initrd | gzip > rootfs/initrd.img

FROM ubuntu:bionic AS grub

RUN apt-get update \
&& apt-get install -y wget build-essential autoconf automake python bison flex

RUN wget -nv -nc ftp://ftp.gnu.org/gnu/grub/grub-2.02.tar.xz \
&& wget -nv -nc https://raw.githubusercontent.com/openwrt/openwrt/master/package/boot/grub2/patches/100-grub_setup_root.patch \
&& tar xf grub-2.02.tar.xz \
&& export CFLAGS=-Wno-error \
&& cd grub-2.02 \
&& patch -p1 < ../100-grub_setup_root.patch \
&& ./autogen.sh \
&& ./configure \
&& make \
&& make install DESTDIR=/grub

FROM ubuntu:bionic

# Copy the kernel and system images "as is",
# extract the compressed initrd and ramdisk images

WORKDIR /android/rootfs
COPY --from=kernel /kernel .
COPY --from=compressed /android/rootfs/ramdisk.img ramdisk.img
COPY --from=compressed /android/rootfs/initrd.img initrd.img

WORKDIR /android/rootfs/system
COPY --from=system / .

# Make the vhd
RUN mkdir data

WORKDIR /android
COPY grub.cfg /android/rootfs/grub/grub.cfg
COPY android-x86.sfdisk .

ENV image_name=android-x86
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
#	-r 1 revision level, default
#	-N 0 number of inodes
#	-m 5 reserved blocks percentage
# label
	-L android \
	$image_name.img \
	-E offset=$((2048 * 512)) \
# blocks count, 1024 * block size (4096) = 4 GB,
# but there's a 1 MB boot sector
	$(( (image_size - 1) * 1024))


COPY --from=grub /grub/ /
ENV GRUB2_MODULES="biosdisk boot chain configfile ext2 linux ls part_msdos reboot serial vga"

COPY grub-early.cfg .

RUN mkdir -p grub2/ \
&& ls -l /usr/local/bin \
&& ls -l /usr/local/lib/grub/i386-pc \
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
&& echo "(hd0) /android/${image_name}.img" > device.map \
&& echo "(hd0) overlay" >> device.map \
&& cat device.map \
&& /usr/local/sbin/grub-bios-setup --help \
&& /usr/local/sbin/grub-bios-setup \
        --device-map=device.map \
        -d grub2/ \
        /android/$image_name.img \
	-r "hd0,msdos1"
