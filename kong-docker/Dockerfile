FROM kong:0.14.1-alpine

WORKDIR /usr/local/kong

RUN apk update && apk add postgresql-client

COPY entrypoint.sh .
COPY plugins/ /usr/local/share/lua/5.1/kong/plugins/

ARG KONG_CUSTOM_PLUGINS
ENV KONG_CUSTOM_PLUGINS=$KONG_CUSTOM_PLUGINS

CMD ["/usr/local/kong/entrypoint.sh"]
