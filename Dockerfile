FROM frolvlad/alpine-glibc
MAINTAINER Danny Krainas <me@danielkrainas.com>

ENV SHEXD_CONFIG_PATH /etc/shexd.default.yml

COPY ./dist /bin/shexd
COPY ./config.default.yml /etc/shexd.default.yml

ENTRYPOINT ["/bin/shexd"]
