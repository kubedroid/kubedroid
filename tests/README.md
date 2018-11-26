# Steps to smoke test the latest images

## With qemu on the host

Make sure you've actually installed qemu. Then, extract the .qcow2 image from the disk
image, and start a local instance of qemu. You can open a VNC connection to connect
to the VM.

The commands below start the VM with and without GPU acceleration.

```
# Libepoxy from source
apt-get -y install git build-essential autoconf autogen libtool pkg-config xutils-dev libgles2-mesa-dev
git clone https://github.com/anholt/libepoxy
./autogen.sh
make -j8
make install

# Tools required to compile virglrenderer
apt-get -y install check libgbm-dev
git clone https://github.com/freedesktop/virglrenderer
cd virglrenderer
CFLAGS='-g -O0' CXXFLAGS='-g -O0' ./autogen.sh --enable-debug=yes --enable-tests
make -j8
make install

# qemu
apt-get install -y libz-dev libglib2.0-dev libplixman-1-dev
wget https://download.qemu.org/qemu-3.0.0.tar.xz
tar xvJf qemu-3.0.0.tar.xz
cd qemu-3.0.0
CFLAGS='-g -O0' CXXFLAGS='-g -O0' ./configure --enable-virglrenderer --enable-vnc --target-list="x86_64-softmmu" --disable-sdl --disable-gtk --enable-system --disable-user --enable-debug
make -j8
make install
```

```
cd ~
docker pull quay.io/quamotion/android-x86-disk:7.1-r2
docker run -v $(pwd):/target --rm quay.io/quamotion/android-x86-disk:7.1-r2 /bin/bash -c "cp /disk/android-x86.qcow2 /target/android-x86-7.1-r2.qcow2"

docker pull quay.io/quamotion/android-x86-disk:8.1-rc2
docker run -v $(pwd):/target --rm quay.io/quamotion/android-x86-disk:8.1-rc2 /bin/bash -c "cp /disk/android-x86.qcow2 /target/android-x86-8.1-rc2.qcow2"

export LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH
qemu-system-x86_64 \
    -enable-kvm \
    -m 2048 -smp 2 -cpu host \
    -device virtio-mouse-pci -device virtio-keyboard-pci \
    -serial mon:stdio \
    -netdev user,id=mynet,hostfwd=tcp::5555-:5555 -device virtio-net-pci,netdev=mynet \
    -vga virtio \
    -display egl-headless -vnc :0 \
    -hda android-x86-7.1-r2.qcow2 \
    -D ~/qemu-logs
```

Or debug:

```
apt-get install -y gdb
gdb --args qemu-system-x86_64     -enable-kvm     -m 2048 -smp 2 -cpu host     -device virtio-mouse-pci -device virtio-keyboard-pci     -serial mon:stdio     -netdev user,id=mynet,hostfwd=tcp::5555-:5555 -device virtio-net-pci,netdev=mynet     -vga virtio     -display egl-headless -vnc 0.0.0.0:5900     -hda android-x86.qcow2 -D ~/qemu-logs
handle SIGUSR1 nostop
run
bt
frame 4
info locals
```
