# Build: docker build -f Dockerfile addons/konzatsu -t konzatsu .
FROM golang:1.27-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /konzatsu .

FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=builder /konzatsu /konzatsu
USER nonroot
EXPOSE 2112
ENTRYPOINT ["/konzatsu"]
