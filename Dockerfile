#Builder
FROM golang:1.18.3-alpine3.16 as builder

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /

COPY ./config.yaml ./app/config.yaml
COPY ./rds-combined-ca-bundle.pem ./app/rds-combined-ca-bundle.pem
COPY ./credentials.json ./app/credentials.json

WORKDIR /app

COPY . .

RUN make engine

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app

WORKDIR /app

EXPOSE 9090

COPY --from=builder /app/engine /app
COPY --from=builder /app/config.yaml /app/config.yaml
COPY --from=builder /app/rds-combined-ca-bundle.pem /app/rds-combined-ca-bundle.pem
COPY --from=builder /app/credentials.json /app/credentials.json

CMD /app/engine