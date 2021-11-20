ARG GO_VERSION

FROM golang:${GO_VERSION}-buster

ARG GORELEASER_VERSION
RUN curl -OL https://github.com/goreleaser/goreleaser/releases/download/v${GORELEASER_VERSION}/goreleaser_Linux_x86_64.tar.gz \
 && ls -ahl --color \
 && tar xvf goreleaser_Linux_x86_64.tar.gz \
 && mv goreleaser /usr/local/bin

ENTRYPOINT ["goreleaser"]
