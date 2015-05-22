FROM golang@1.4.2

ENV http_proxy http://proxy.***REMOVED***:8080
ENV https_proxy http://proxy.***REMOVED***:8080
ENV no_proy sap.corp,localhost,127.0.0.1

RUN mkdir -p /gonative
RUN go get github.com/mitchellh/gox
RUN go get github.com/inconshreveable/gonative
RUN cd /gonative && gonative build
ENV PATH /gonative/go/bin:$PATH

ENV GOPATH /gonative/src/gitHub.***REMOVED***/monsoon/arc/Godeps/_workspace:/gonative 
WORKDIR /gonative/src/gitHub.***REMOVED***/monsoon/arc

