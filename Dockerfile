FROM buildpack-deps:wheezy-scm

# gcc for cgo
RUN apt-get update && apt-get install -y \
		gcc libc6-dev make \
		--no-install-recommends \
	&& rm -rf /var/lib/apt/lists/*

RUN curl http://aia.pki.co.sap.com/aia/SAP%20Global%20Root%20CA.crt | \
    tr -d '\r' > /usr/local/share/ca-certificates/SAP_Global_Root_CA.crt && \
    update-ca-certificates

COPY gonative_linux /usr/bin/gonative

RUN gonative build -version 1.4.3 -target=/usr/local/go -platforms "linux_amd64 windows_amd64"

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

RUN go get github.com/mitchellh/gox
RUN go get bitbucket.org/liamstask/goose/cmd/goose
RUN go get github.com/mjibson/esc
RUN go get github.com/blynn/nex
RUN go get github.com/constabulary/gb/cmd/gb
RUN go get golang.org/x/tools/cmd/goimports
