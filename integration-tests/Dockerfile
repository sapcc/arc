FROM progrium/busybox

RUN http_proxy=http://proxy.***REMOVED***:8080 opkg-install ca-certificates && \
    wget -qO- http://aia.pki.co.sap.com/aia/SAP%20Global%20Root%20CA.crt | tr -d '\r' > /etc/ssl/certs/SAP_Global_Root.crt

# satisfy go crypto/x509
RUN for cert in `ls -1 /etc/ssl/certs/*.crt | grep -v /etc/ssl/certs/ca-certificates.crt`; \
      do cat "$cert" >> /etc/ssl/certs/ca-certificates.crt; \
    done

COPY env.sh /env.sh
ENTRYPOINT ["/env.sh"]

COPY bin/ /opt/bin/
ENV PATH /opt/bin:$PATH

CMD ["/bin/sh"]
