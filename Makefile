IMAGE := $(or ${IMAGE}, localhost/test:latest)
PWD = $(shell pwd)
LOG_LEVEL := $(or ${LOG_LEVEL}, info)
ROOT_DIR = $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

build:
	podman build -f Dockerfile . -t $(IMAGE)
