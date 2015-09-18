FROM progrium/busybox

RUN http_proxy=http://proxy.***REMOVED***:8080 opkg-install bash curl bind-dig

COPY env.sh /env.sh
ENTRYPOINT ["/env.sh"]

ADD bin/update-site /usr/bin/update-site

CMD ["update-site", "local"]

