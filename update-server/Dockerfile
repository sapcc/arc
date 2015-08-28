FROM alpine@3.2

ADD bin/update-site /update-site

ENV ARTIFACTS_PATH=/artifacts \
    BIND_ADDRESS=0.0.0.0:8080

CMD ["/update-site"]


