FROM quamotion/android-x86-tools

ENV ISO_URL=https://osdn.net/dl/android-x86/android-x86_64-7.1-r2.iso
ENV ISO_FILE=android-x86_64-7.1-r2.iso

WORKDIR /src

RUN wget $ISO_URL -O $ISO_FILE

RUN isoinfo -i $ISO_FILE -x /system.sfs > system.sfs \
&& isoinfo -i $ISO_FILE -x /initrd.img > initrd.img \
&& isoinfo -i $ISO_FILE -x /ramdisk.img > ramdisk.img \
&& isoinfo -i $ISO_FILE -x /kernel > kernel \
&& unsquashfs -f -d . system.sfs system.img

WORKDIR /android/
RUN cp /src/kernel . \
&& mkdir initrd ramdisk \
&& (cd initrd && zcat /src/initrd.img | cpio -idv) \
&& (cd ramdisk && zcat /src/ramdisk.img | cpio -idv) \
&& 7z x /src/system.img -osystem/ \
&& rm -rf system/'[SYS]' \
&& rm -rf system/lost+found

