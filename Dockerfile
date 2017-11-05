FROM haproxy:1.7

ENV HAPROXY_USER haproxy

RUN groupadd --system ${HAPROXY_USER} \
    && useradd --system --gid ${HAPROXY_USER} ${HAPROXY_USER} \
    && mkdir -p /var/lib/${HAPROXY_USER} \
    && chown -R ${HAPROXY}:${HAPROXY_USER} /var/lib/${HAPROXY_USER}

COPY haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg

CMD ["haproxy", "-f", "/usr/local/etc/haproxy/haproxy.cfg"]

