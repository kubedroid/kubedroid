docker: android-x86_64-7.1-r2.iso
	sudo docker build . -t quay.io/quamotion/android-x86-base:7.1-r2

android-x86_64-7.1-r2.iso:
	wget -nc https://osdn.net/dl/android-x86/android-x86_64-7.1-r2.iso -O android-x86_64-7.1-r2.iso

run:
	sudo docker run --rm -it quay.io/quamotion/android-x86-base:7.1-r2 /bin/bash
