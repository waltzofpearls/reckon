ARG GO_VERSION

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
COPY --from=osxcross /osxcross/target/ /usr/local/osxcross/
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
        gcc-arm-linux-gnueabihf \
        g++-arm-linux-gnueabihf \
        gcc-aarch64-linux-gnu \
        g++-aarch64-linux-gnu \
        clang \
 && rm -rf /var/lib/apt/lists/*
ARG GORELEASER_VERSION
RUN curl -OL https://github.com/goreleaser/goreleaser/releases/download/v${GORELEASER_VERSION}/goreleaser_Linux_x86_64.tar.gz \
 && ls -ahl --color \
 && tar xvf goreleaser_Linux_x86_64.tar.gz \
 && mv goreleaser /usr/local/bin

ENTRYPOINT ["goreleaser"]
