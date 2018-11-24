# Steps to smoke test the latest images

## With qemu on the host

Make sure you've actually installed qemu. Then, extract the .qcow2 image from the disk
image, and start a local instance of qemu. You can open a VNC connection to connect
to the VM.

The commands below start the VM with and without GPU acceleration.

```
apt-get -y install git libglib2.0-dev libfdt-dev libpixman-1-dev zlib1g-dev build-essential autoconf autogen libtool libdrm-dev libgbm-dev libepoxy-dev libx11-dev libegl1-mesa-dev xorg-dev

wget https://github.com/freedesktop/virglrenderer/archive/virglrenderer-0.7.0.tar.gz
tar xvzf virglrenderer-0.7.0.tar.gz
mv virglrenderer-virglrenderer-0.7.0 virglrenderer-0.7.0
cd virglrenderer-0.7.0
./autogen.sh
make
make install

wget https://download.qemu.org/qemu-3.0.0.tar.xz
tar xvJf qemu-3.0.0.tar.xz
cd qemu-3.0.0
./configure --enable-virglrenderer --enable-vnc --target-list="x86_64-softmmu" --disable-sdl --disable-gtk --enable-system --disable-user
make -j8
```

```
docker run -v $(pwd):/target --rm quay.io/quamotion/android-x86:8.1-rc2 /bin/bash -c "cp /disk/* /target/"

qemu-system-x86_64 \
    -enable-kvm \
    -m 2048 -smp 2 -cpu host \
    -device virtio-mouse-pci -device virtio-keyboard-pci \
    -serial mon:stdio \
    -netdev user,id=mynet,hostfwd=tcp::5555-:5555 -device virtio-net-pci,netdev=mynet \
    -vga virtio \
    -display egl-headless -vnc :4444 \
    -hda android-x86.qcow2

```
