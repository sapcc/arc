FROM keppel.eu-de-1.cloud.sap/ccloud-dockerhub-mirror/library/alpine:3.4
LABEL source_repository="https://github.com/sapcc/arc"

RUN apk add --no-cache curl jq
RUN curl -fL http://aia.pki.co.sap.com/aia/SAP%20Global%20Root%20CA.crt | tr -d '\r' > /usr/local/share/ca-certificates/SAP_Global_Root_CA.crt \
      && update-ca-certificates

COPY bin/ /opt/bin/
COPY get-token run_tests.sh /opt/bin/
ENV PATH /opt/bin:$PATH

CMD ["run_tests.sh"]
