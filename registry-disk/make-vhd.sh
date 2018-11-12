#!/bin/bash
image_name=android-x86

lofile=`losetup -f`

losetup -P $lofile $image_name.img

mkdir -p /mnt/$image_name
mount ${lofile}p1 /mnt/$image_name

echo "(hd0) $lofile" > $image_name.map
grub-install --no-floppy --grub-mkdevicemap=$image_name.map --modules="part_msdos" --boot-directory=/mnt/$image_name $lofile

losetup -d $lofile

qemu-img convert -f raw -O vpc $image_name.img $image_name.vhd
