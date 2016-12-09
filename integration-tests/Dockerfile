FROM alpine:3.4 

RUN apk add --no-cache curl jq
RUN curl -fL http://aia.pki.co.sap.com/aia/SAP%20Global%20Root%20CA.crt | tr -d '\r' > /usr/local/share/ca-certificates/SAP_Global_Root_CA.crt \
      && update-ca-certificates \
			&& curl https://gitHub.***REMOVED***

COPY bin/ /opt/bin/
COPY get-token run_tests.sh /opt/bin/
ENV PATH /opt/bin:$PATH

CMD ["run_tests.sh"]
