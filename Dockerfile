FROM registry.access.redhat.com/ubi9/go-toolset:1.21 AS golang

USER 0

## Build
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=1 GOFLAGS="" GO111MODULE=on go build -o /test main.go

## Runtime
FROM quay.io/centos/centos:stream9

ARG release=main
ARG version=latest
ARG DATA_DIR=/data
ARG UID=1001
ARG GID=0
RUN mkdir $DATA_DIR && chmod 775 $DATA_DIR && chown $UID:$GID /data
VOLUME $DATA_DIR
ENV DATA_DIR=$DATA_DIR

RUN dnf install -y epel-release
RUN dnf install -y p7zip-plugins squashfs-tools && dnf clean all
COPY data/rootfs.img /tmp

USER $UID:$GID
COPY --from=golang /test /test
CMD ["/test"]
