FROM progrium/busybox


COPY env.sh /env.sh
ENTRYPOINT ["/env.sh"]

ADD bin/update-site /usr/bin/update-site

CMD ["update-site", "local"]

