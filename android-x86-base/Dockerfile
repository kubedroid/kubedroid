FROM quay.io/quamotion/android-x86-tools

ENV ISO_URL=https://osdn.net/dl/android-x86/android-x86_64-7.1-r2.iso
ENV ISO_FILE=android-x86_64-7.1-r2.iso

WORKDIR /android

# Download the ISO file
RUN wget $ISO_URL -O $ISO_FILE \

# Extract the root file system
&& isoinfo -i $ISO_FILE -x /system.sfs > system.sfs \
&& unsquashfs -f -d . system.sfs system.img \
&& 7z x system.img -osystem/ \
&& rm -rf system/'[SYS]' system/lost+found \
&& rm system.sfs system.img \

# Extract the initrd image
&& isoinfo -i $ISO_FILE -x /initrd.img > initrd.img \
&& mkdir initrd \
&& (cd initrd && zcat ../initrd.img | cpio -idv) \
&& rm initrd.img \

# Extract the ramdisk image
&& isoinfo -i $ISO_FILE -x /ramdisk.img > ramdisk.img \
&& mkdir ramdisk \
&& (cd ramdisk && zcat ../ramdisk.img | cpio -idv) \
&& rm ramdisk.img \

# Extract the kernel
&& isoinfo -i $ISO_FILE -x /kernel > kernel \

# Remove the ISO file
&& rm $ISO_FILE

