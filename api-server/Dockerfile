FROM progrium/busybox

ADD https://gitHub.***REMOVED***/userID/goose/releases/download/0.1/goose_linux_amd64_static /usr/bin/goose
RUN chmod +x /usr/bin/goose

RUN http_proxy=http://proxy.***REMOVED***:8080 opkg-install ca-certificates && \
    wget -qO- http://aia.pki.co.sap.com/aia/SAP%20Global%20Root%20CA.crt | tr -d '\r' > /etc/ssl/certs/SAP_Global_Root.crt

# satisfy go crypto/x509
RUN for cert in `ls -1 /etc/ssl/certs/*.crt | grep -v /etc/ssl/certs/ca-certificates.crt`; \
      do cat "$cert" >> /etc/ssl/certs/ca-certificates.crt; \
    done

COPY env.sh /env.sh
ENTRYPOINT ["/env.sh"]

ENV ARC_ENV=development \
    DBHOST=postgres \
    DBUSER=postgres

WORKDIR /opt/api
COPY db/ db/
COPY bin/gencert /usr/bin/gencert
COPY bin/api-server /usr/bin/api-server

CMD goose -env=$ARC_ENV status && goose -env=$ARC_ENV up && exec api-server
