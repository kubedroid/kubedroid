FROM ubuntu:bionic AS grub

RUN apt-get update \
&& apt-get install -y wget build-essential autoconf automake python bison flex \
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

FROM ubuntu:bionic

RUN apt-get update \
&& apt-get install -y qemu-utils \
&& rm -rf /var/lib/apt/lists*

COPY --from=grub /grub/ /
