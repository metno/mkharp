FROM ubuntu:18.04

ARG GOLANG_VERSION=1.17.1
ARG GOLANG_SHA256=dab7d9c34361dc21ec237d584590d72500652e7c909bf082758fb63064fca0ef

RUN apt-get update && apt-get -y install git curl build-essential

RUN \
    curl -L -O https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    sha256sum go${GOLANG_VERSION}.linux-amd64.tar.gz | grep -q ${GOLANG_VERSION} && \
    tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    rm go${GOLANG_VERSION}.linux-amd64.tar.gz

ENV GOPATH=/go
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

ENV DEBIAN_FRONTEND=noninteractive

COPY go.mod /src/go.mod
COPY go.sum /src/go.sum
WORKDIR /src
RUN go mod download

COPY internal /src/internal
COPY main.go /src/

RUN go test ./...
RUN go install