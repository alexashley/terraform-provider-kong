FROM kong:0.14.1-alpine

WORKDIR /usr/local/kong

RUN apk update && apk add postgresql-client

COPY . .

CMD ["/usr/local/kong/entrypoint.sh"]
