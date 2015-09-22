FROM progrium/busybox

RUN http_proxy=http://proxy.***REMOVED***:8080 opkg-install bash curl bind-dig

RUN opkg-install ca-certificates
ADD SAP_Global_Root_CA.crt /etc/ssl/certs/SAP_Global_Root_CA.crt
# satisfy go crypto/x509
RUN for cert in `ls -1 /etc/ssl/certs/*.crt | grep -v /etc/ssl/certs/ca-certificates.crt`; \
      do cat "$cert" >> /etc/ssl/certs/ca-certificates.crt; \
    done

COPY env.sh /env.sh
ENTRYPOINT ["/env.sh"]

ADD bin/update-site /usr/bin/update-site

CMD ["update-site", "local"]

