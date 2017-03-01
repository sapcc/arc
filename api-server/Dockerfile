FROM alpine:3.4 

RUN apk add --no-cache curl
RUN curl -fL http://aia.pki.co.sap.com/aia/SAP%20Global%20Root%20CA.crt | tr -d '\r' > /usr/local/share/ca-certificates/SAP_Global_Root_CA.crt \
      && update-ca-certificates
RUN curl -fLo /usr/bin/goose https://gitHub.***REMOVED***/userID/goose/releases/download/0.1/goose_linux_amd64_static \
      && chmod +x /usr/bin/goose

ENV ARC_ENV=development \
    DBHOST=postgres \
    DBUSER=postgres

WORKDIR /opt/api
COPY entrypoint.sh /entrypoint.sh
COPY db/ db/
COPY etc/ etc/
COPY bin/api-server /usr/bin/api-server
COPY bin/janitor /usr/bin/janitor

ENTRYPOINT ["/entrypoint.sh"]
CMD ["api-server"]
