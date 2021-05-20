FROM golang:alpine as builder
WORKDIR /app
RUN apk add --update --no-cache ca-certificates && update-ca-certificates
COPY main.go go.* /app/
RUN CGO_ENABLED=0 GOOS=linux go build

FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/discord_time /bin/discord_time
ENTRYPOINT [ "/bin/discord_time" ]
