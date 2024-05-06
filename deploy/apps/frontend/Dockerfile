FROM golang:1.22 as builder

COPY [".", "/build"]
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /build/frontend /build/server/cmd/frontend

FROM alpine:3.19

ENV CONFIG_PATH /cfg/local.yaml

COPY --from=builder /build/server/config /cfg
COPY --from=builder /build/frontend /usr/local/bin/frontend

RUN chmod +x /usr/local/bin/frontend

CMD frontend