ARG BUILD_FROM=alpine:3.22

FROM golang:1.27-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /konzatsu .

FROM $BUILD_FROM
COPY --from=builder /konzatsu /usr/bin/konzatsu
COPY run.sh /
RUN chmod +x /run.sh
ENTRYPOINT ["/run.sh"]
