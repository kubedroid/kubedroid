# This image starts from kubevirt/virt-launcher:v0.10.0, patches it, and
# is intended to then replace it.
# To make sure we always start from the original upstream copy, we start
# from the SHA hash instead of the image tag (which would be replaced)
#
# The main goal is to get an up-to-date version of libvirt and qemu,
# so that we can hardware-accelerate Android-x86
FROM golang:1.11 AS build

WORKDIR /go/src/virt-launcher

COPY *.go .
RUN go build

FROM 9b33dab9e7d75ecbbbb4c599dd9c0e1ab9f6f126fe605989bc86d8a7e80de9fc

ARG LIBVIRT_VERSION=4.9.0
ARG QEMU_VERSION=3.0.0
ARG LIBUSB_VERSION=1.0.22

RUN dnf install -y dnf-plugins-core \
&& dnf copr enable -y @virtmaint-sig/virt-preview \
&& dnf install -y \
     libvirt-daemon-kvm-$LIBVIRT_VERSION \
     libvirt-client-$LIBVIRT_VERSION \
     qemu-kvm-$QEMU_VERSION \
     libusbx-$LIBUSB_VERSION \
      mesa-dri-drivers \
&& dnf clean all

RUN mv /usr/bin/virt-launcher /usr/bin/upstream-virt-launcher
COPY --from=build /go/src/virt-launcher/virt-launcher /usr/bin/
