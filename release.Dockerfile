ARG PYTHON_VERSION
ARG GO_VERSION

FROM amd64/python:${PYTHON_VERSION}-slim-buster AS amd64
FROM arm32v7/python:${PYTHON_VERSION}-slim-buster AS armhf
FROM arm64v8/python:${PYTHON_VERSION}-slim-buster AS arm64

FROM debian:buster-slim AS osxcross
ARG OSX_SDK_VERSION
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
        git \
        ca-certificates \
        curl \
        make \
        python3 \
        clang \
        cmake \
        patch \
        libxml2-dev \
        libssl-dev \
        zlib1g-dev \
        xz-utils \
 && rm -rf /var/lib/apt/lists/*
RUN git clone https://github.com/tpoechtrager/osxcross.git \
 && cd osxcross/tarballs \
 && curl -OL https://github.com/phracker/MacOSX-SDKs/releases/download/11.3/MacOSX${OSX_SDK_VERSION}.sdk.tar.xz
RUN cd osxcross \
 && UNATTENDED=yes ./build.sh

FROM golang:${GO_VERSION}-buster
ENV AMD64=/usr/local/python_amd64
ENV ARMHF=/usr/local/python_armhf
ENV ARM64=/usr/local/python_arm64
# copy files for amd64
COPY --from=amd64 /usr/local/include/python3.7m ${AMD64}/include/python3.7m
COPY --from=amd64 /usr/local/lib/libpython3*.so* ${AMD64}/lib/
COPY --from=amd64 /usr/local/lib/pkgconfig/ ${AMD64}/lib/pkgconfig/
COPY --from=amd64 /usr/local/lib/python3.7/ ${AMD64}/lib/python3.7/
# copy files for armhf (arm32v7)
COPY --from=armhf /usr/local/include/python3.7m ${ARMHF}/include/python3.7m
COPY --from=armhf /usr/local/lib/libpython3*.so* ${ARMHF}/lib/
COPY --from=armhf /usr/local/lib/pkgconfig/ ${ARMHF}/lib/pkgconfig/
COPY --from=armhf /usr/local/lib/python3.7/ ${ARMHF}/lib/python3.7/
# copy files for arm64 (arm64v8)
COPY --from=arm64 /usr/local/include/python3.7m ${ARM64}/include/python3.7m
COPY --from=arm64 /usr/local/lib/libpython3*.so* ${ARM64}/lib/
COPY --from=arm64 /usr/local/lib/pkgconfig/ ${ARM64}/lib/pkgconfig/
COPY --from=arm64 /usr/local/lib/python3.7/ ${ARM64}/lib/python3.7/
RUN cd ${AMD64}/lib \
 && rm -f libpython3.7m.so \
 && ln -s libpython3.7m.so.1.0 libpython3.7m.so \
 && sed -i "s|prefix=/usr/local|prefix=${AMD64}|" pkgconfig/python3.pc \
 && cd ${ARMHF}/lib \
 && rm -f libpython3.7m.so \
 && ln -s libpython3.7m.so.1.0 libpython3.7m.so \
 && sed -i "s|prefix=/usr/local|prefix=${ARMHF}|" pkgconfig/python3.pc \
 && cd ${ARM64}/lib \
 && rm -f libpython3.7m.so \
 && ln -s libpython3.7m.so.1.0 libpython3.7m.so \
 && sed -i "s|prefix=/usr/local|prefix=${ARM64}|" pkgconfig/python3.pc
# copy files from osxcross
COPY --from=osxcross /osxcross/target/ /usr/local/osxcross/
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
        gcc-arm-linux-gnueabihf \
        g++-arm-linux-gnueabihf \
        gcc-aarch64-linux-gnu \
        g++-aarch64-linux-gnu \
 && rm -rf /var/lib/apt/lists/*
ARG GORELEASER_VERSION
RUN curl -OL https://github.com/goreleaser/goreleaser/releases/download/v${GORELEASER_VERSION}/goreleaser_Linux_x86_64.tar.gz \
 && ls -ahl --color \
 && tar xvf goreleaser_Linux_x86_64.tar.gz \
 && mv goreleaser /usr/local/bin

ENTRYPOINT ["goreleaser"]
