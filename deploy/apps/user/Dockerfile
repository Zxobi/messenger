FROM golang:1.22 as builder

COPY [".", "/build"]
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /build/user /build/server/cmd/user

FROM alpine:3.19

ENV CONFIG_PATH /cfg/local.yaml

COPY --from=builder /build/server/config /cfg
COPY --from=builder /build/user /usr/local/bin/user

RUN chmod +x /usr/local/bin/user

CMD user