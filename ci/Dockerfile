FROM ubuntu:14.04
MAINTAINER Ferran Rodenas <frodenas@gmail.com>

# Install dependencies
RUN apt-get update && apt-get -y install build-essential curl git-core

# Install Golang
ENV GOLANG_VERSION 1.4.2
RUN curl -sSL https://storage.googleapis.com/golang/go$GOLANG_VERSION.linux-amd64.tar.gz | tar -v -C /usr/local -xz
ENV GOROOT /usr/local/go
ENV PATH $PATH:$GOROOT/bin
