#!/bin/bash

cd ~
git clone https://github.com/qmfrederik/kube-virt-droid-images
cd ~/kube-virt-droid-images/android-x86-tools
make

cd ~/kube-virt-droid-images/android-x86-base
make

cd ~/kube-virt-droid-images/android-x86-disk
make

cd ~/kube-virt-droid-images/android-x86-hook
make

cd ~/kube-virt-droid-images/android-x86-launcher-image
sudo docker pull docker.io/kubevirt/virt-launcher:v0.10.0
make
