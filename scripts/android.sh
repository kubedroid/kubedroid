#!/bin/bash
sudo docker pull quay.io/quamotion/android-x86-disk:7.1-r2
sudo docker pull quay.io/quamotion/android-x86-disk:8.1-rc2
sudo docker pull quay.io/quamotion/android-x86-hook:master
sudo docker pull quay.io/quamotion/android-x86-launcher:master

# This will force kubevirt to use the Android launcher. Use with care
sudo docker tag quay.io/quamotion/android-x86-launcher:master docker.io/kubevirt/virt-launcher:v0.10.0
