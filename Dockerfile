FROM golang:latest as builder
WORKDIR /application
COPY . .
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o fbackup .

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=builder /application/fbackup /usr/bin/
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["fbackup"]