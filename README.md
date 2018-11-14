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

## References

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

``
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
