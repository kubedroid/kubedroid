FROM ubuntu:bionic AS grub

RUN apt-get update \
&& apt-get install -y wget build-essential autoconf automake python bison flex nano patch \
&& rm -rf /var/lib/apt/lists/*

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

FROM ubuntu:bionic AS mkbootfs
RUN apt-get update \
&& apt-get install -y build-essential android-libcutils-dev git \
&& rm -rf /var/lib/apt/lists

RUN git clone --depth=1 git://git.osdn.net/gitroot/android-x86/system-core.git -b nougat-x86 bootimg-tools \
&& cd bootimg-tools/cpio \
&& gcc mkbootfs.c -o mkbootfs -I../include -lcutils -L/usr/lib/x86_64-linux-gnu/android/ \
&& cp mkbootfs /usr/local/bin/

FROM ubuntu:bionic AS extfstools
RUN apt-get update \
&& apt-get install -y build-essential git clang libboost-dev libstdc++-5-dev \
&& rm -rf /var/lib/apt/lists

RUN git clone --depth=1 https://github.com/qmfrederik/extfstools \
&& cd extfstools \
&& make \
&& cp ext2rd /usr/local/bin \
&& cp ext2dump /usr/local/bin

FROM ubuntu:bionic

COPY --from=grub /grub /
COPY --from=mkbootfs /usr/local/bin/mkbootfs /usr/local/bin/
COPY --from=extfstools /usr/local/bin/ext2rd /usr/local/bin/
COPY --from=extfstools /usr/local/bin/ext2dump /usr/local/bin/

RUN apt-get update \
&& apt-get install -y qemu-utils android-libcutils-dev wget genisoimage squashfs-tools cpio e2tools \
&& rm -rf /var/lib/apt/lists*

ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/lib/x86_64-linux-gnu/android
