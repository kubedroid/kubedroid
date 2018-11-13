FROM ubuntu:bionic AS build
RUN apt-get update \
&& apt-get install -y build-essential android-libcutils-dev git \
&& rm -rf /var/lib/apt/lists

RUN git clone --depth=1 git://git.osdn.net/gitroot/android-x86/system-core.git -b nougat-x86 bootimg-tools \
&& cd bootimg-tools/cpio \
&& gcc mkbootfs.c -o mkbootfs -I../include -lcutils -L/usr/lib/x86_64-linux-gnu/android/ \
&& cp mkbootfs /usr/local/bin/

FROM ubuntu:bionic
RUN apt-get update \
&& apt-get install -y android-libcutils-dev \
&& rm -rf /var/lib/apt/lists

COPY --from=build /usr/local/bin/mkbootfs /usr/local/bin/
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/lib/x86_64-linux-gnu/android
