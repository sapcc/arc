FROM alpine@3.2

ADD bin/api-server /api-server

ENV ARTIFACTS_PATH=/artifacts \
    BIND_ADDRESS=0.0.0.0:8080

CMD ["/api-server"]


