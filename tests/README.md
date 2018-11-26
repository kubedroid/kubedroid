# Steps to smoke test the latest images

## With qemu on the host

Make sure you've actually installed qemu. Then, extract the .qcow2 image from the disk
image, and start a local instance of qemu. You can open a VNC connection to connect
to the VM.

The commands below start the VM with and without GPU acceleration.

```
# Libepoxy from source
cd ~
apt-get update
apt-get -y install git build-essential autoconf autogen libtool pkg-config xutils-dev libgles2-mesa-dev
wget https://github.com/anholt/libepoxy/releases/download/1.5.3/libepoxy-1.5.3.tar.xz
tar xvJf libepoxy-1.5.3.tar.xz
cd libepoxy-1.5.3
./autogen.sh
make -j8
make install

# Tools required to compile virglrenderer
cd ~
apt-get -y install check libgbm-dev
wget https://github.com/freedesktop/virglrenderer/archive/virglrenderer-0.7.0.tar.gz
tar xvzf virglrenderer-0.7.0.tar.gz
mv virglrenderer-virglrenderer-0.7.0/ virglrenderer-0.7.0/
cd virglrenderer-0.7.0
# CFLAGS='-g -O0' CXXFLAGS='-g -O0' ./autogen.sh --enable-debug=yes --enable-tests
./autogen.sh --enable-tests
make -j8
make install

# qemu
cd ~
apt-get install -y libz-dev libglib2.0-dev libpixman-1-dev
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

Or, from qemu/master with the following patch to force surfaceless rendering:

```patch
diff --git a/hw/display/virtio-gpu-3d.c b/hw/display/virtio-gpu-3d.c
index 55d76405a9..7482e37063 100644
--- a/hw/display/virtio-gpu-3d.c
+++ b/hw/display/virtio-gpu-3d.c
@@ -624,7 +624,7 @@ int virtio_gpu_virgl_init(VirtIOGPU *g)
 {
     int ret;

-    ret = virgl_renderer_init(g, 0, &virtio_gpu_3d_cbs);
+    ret = virgl_renderer_init(g, VIRGL_RENDERER_USE_SURFACELESS | VIRGL_RENDERER_USE_EGL, &virtio_gpu_3d_cbs);
     if (ret != 0) {
         return ret;
     }
```

To try with an Ubuntu guest (this works; the default resolution appears to be 1024x768 at bith depth 24):

```
wget http://releases.ubuntu.com/18.04/ubuntu-18.04.1-desktop-amd64.iso
qemu-img create -f qcow2 ubuntu-18.04-desktop-amd64.img.qcow2 16G
qemu-system-x86_64 \
  -cdrom ubuntu-18.04.1-desktop-amd64.iso \
  -drive file=ubuntu-18.04-desktop-amd64.img.qcow2,format=qcow2 \
  -enable-kvm \
  -m 2G \
  -smp 2 \
  -vga virtio \
  -display egl-headless
```
