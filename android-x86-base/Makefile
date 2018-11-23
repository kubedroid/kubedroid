docker: android-x86_64-8.1-rc2.iso
	sudo docker build . -t quay.io/quamotion/android-x86-base:8.1-rc2

android-x86_64-8.1-rc2.iso:
	wget -nc https://osdn.net/dl/android-x86/android-x86_64-8.1-rc2.iso -O android-x86_64-8.1-rc2.iso

run:
	sudo docker run --rm -it quay.io/quamotion/android-x86-base:8.1-rc2 /bin/bash
