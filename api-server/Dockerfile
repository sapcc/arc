FROM progrium/busybox

ADD https://gitHub.***REMOVED***/userID/goose/releases/download/0.1/goose_linux_amd64_static /usr/bin/goose
RUN chmod +x /usr/bin/goose

COPY env.sh /env.sh
ENTRYPOINT ["/env.sh"]

ENV ARC_ENV=development \
    DBHOST=postgres \
    DBUSER=postgres

WORKDIR /opt/api
COPY db/ db/
COPY bin/api-server /usr/bin/api-server

CMD goose -env=$ARC_ENV status && goose -env=$ARC_ENV up && exec api-server
