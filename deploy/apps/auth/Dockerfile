FROM golang:1.22 as builder

COPY [".", "/build"]
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /build/auth /build/server/cmd/auth

FROM alpine:3.19

ENV CONFIG_PATH /cfg/local.yaml

COPY --from=builder /build/server/config /cfg
COPY --from=builder /build/auth /usr/local/bin/auth

RUN chmod +x /usr/local/bin/auth

CMD auth