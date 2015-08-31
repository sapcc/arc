FROM buildpack-deps:wheezy-scm

ENV http_proxy=http://proxy.***REMOVED***:8080 \
    https_proxy=http://proxy.***REMOVED***:8080 \
    no_proy***REMOVED***,localhost,127.0.0.1

# gcc for cgo
RUN apt-get update && apt-get install -y \
		gcc libc6-dev make \
		--no-install-recommends \
	&& rm -rf /var/lib/apt/lists/*

COPY gonative_linux /usr/bin/gonative

RUN gonative build -version 1.4.2 -target=/usr/src/go -platforms "linux_amd64 windows_amd64"

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/src/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

RUN go get github.com/mitchellh/gox
RUN go get bitbucket.org/liamstask/goose/cmd/goose 
RUN go get github.com/mjibson/esc 
RUN go get github.com/blynn/nex 
