FROM golang:latest as builder
WORKDIR /application
COPY . .
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o fbackup .

FROM scratch
COPY --from=builder /application/fbackup /usr/bin/
ENTRYPOINT ["fbackup"]