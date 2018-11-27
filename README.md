# Android-x86 images for KubeVirt

This repository contains scripts which are used to created Android-x86 images
for use with KubeVirt.

The output is a docker image, quay.io/quamotion/android-x86-disk:7.1-r2, which
you can use as a registryVolume in KubeVirt when provisioning your VM.

To create an Android VM, you can run the following command:

```
kubectl create -f android-x86.yml
```

to provision a VM running Android-x86.

## Console, VNC and adb access

You can access your Android machine using the text-based console,
VNC and adb.

To connect over a serial console, type:

```
virtctl console android-x86
```

You may need to hit ENTER to get a command prompt. Please note
that the console is not accessible when the device is in sleep
mode.

To connect over VNC, type:

```
virtctl vnc android-x86
```

This will get you VNC access to your Android-x86 VM.

To connect over ADB, first get the IP address for your VM:

```
kubectl get vmi
```

Then, connect to port 5555:

```
adb connect <ip>:5555
```

Or, in one command:

```
adb connect `kubectl  get -o=jsonpath='{.status.interfaces[0].ipAddress}' vmi android-x86`:555
```

## Initializing your VM
`adb` is enabled by default, so you can `adb connect [VM IP]:5555` to get a connection to your VM.

You may want to configure your VM over adb. For example, you may want to:

- Enable *unknown sources*: `adb shell settings put secure install_non_market_apps 1`
- 

## References

### Hardware acceleration

You may like your Android-x86 VM better if it has hardware acceleration. [You can use Android-x86 with
virtio (virglrenderer)](https://groups.google.com/forum/#!msg/android-x86/enPcst6oQ_w/8Etr0aEZAAAJ).

With vanilla QEMU on a desktop, the command line would be:

```
qemu-system-x86_64 \
    -enable-kvm \
    -m 2048 -smp 2 -cpu host \
    -device virtio-mouse-pci -device virtio-keyboard-pci \
    -serial mon:stdio \
    -boot menu=on \
    -netdev user,id=mynet,hostfwd=tcp::5555-:5555 -device virtio-net-pci,netdev=mynet \
    -vga virtio -display sdl,gl=on $@ \
    -cdrom ${ANDROID_IMAGE_PATH} \
    -hda android.img
```

But:

> `display sdl` means a local (to the node) display will be spawned using SDL

and

> `-spice gl=on,...` works only locally (qemu and spice client must run on the same machine).
> `-display egl-headless -spice gl=off,...` works remotely.  Not very efficient, it'll effectively do ReadPixels on the rendered framebuffer and send them over spice like non-gl display updates.  But it might still be better than using llvmpipe in the android guest.
> `-display egl-headless -vnc ...` works too.


Net, for this to work with KubeVirt, [we need three things](https://groups.google.com/d/msg/kubevirt-dev/7xYZQtILpJM/KtTqLnO9AAAJ):

1. The latest software:
   a. [libvirt 4.6.0 or newer](https://github.com/libvirt/libvirt/commit/d8266ebe1615c4b043db6b8d486465722cdd0ef8) (you'll need KubeVirt master)
   b. [qemu 2.10 or newer](https://patchwork.kernel.org/patch/10465793/)
2. Configure QEMU to use `/dev/dri/renderD*` 
3. Launch QEMU with the correct arguments, so that it uses virtio

The plan to get this working is to:
1. Use a KubeVirt hook to update the domain definition, this is implemented in the android-x86-hook container.
2. Use the KubeVirt master which has a sufficiently recent version of QEMU and libvirt
3. Figure out how to pass the `/dev/dri/render*` device to QEMU

The [libvirt user manual](https://libvirt.org/formatdomain.html#elementsGraphics) has more information
about the different graphic options available.

### Making an Android-x86 disk image

You can create a custom .vhd image which you can use to boot Android-x86 in a VM, using the
Android-x86 installer .iso file, and the custom kernel. Make sure to use ext4.

References:

- https://forum.xda-developers.com/android/general/guide-triple-boot-android-ubuntu-t3092913
- http://my-zhang.github.io/blog/2014/06/28/make-bootable-linux-disk-image-with-grub2/
- https://superuser.com/questions/130955/how-to-install-grub-into-an-img-file

You'll need a patched copy of grub2 to make it work properly with images without having
to mount them.

### Rebuilding the ramdisk
Unpacking, modifying and repacking the ramdisk is as simple as:

```
mkdir ramdisk
cd ramdisk
gzip -dc ../ramdisk.img | cpio -i
```
and

```
cd ..
mkbootfs ./ramdisk | gzip > ramdisk_new.gz
mv ramdisk_new.gz ramdisk_new.img
once you have mkbootfs.
```

You'll need to build mkbootfs from the Android source, and you need libcutils.

In short:

```
cd system/core/cpio
gcc mkbootfs.c -o mkbootfs -I../include -lcutils -L/usr/lib/android/
export LD_LIBRARY_PATH=/usr/lib/android
./mkbootfs
```

### Running Android in Docker

Not quite there yet, but see:
- https://github.com/Rudloff/termux-docker-image/
- https://github.com/qmfrederik/docker-android-x86

## Project Sponsor

This project is sponsored by [Quamotion](http://quamotion.mobi).
