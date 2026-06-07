# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY main.go index.html ./

ARG TARGETOS=linux
ARG TARGETARCH=amd64

ENV CGO_ENABLED=0
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -o weather .

FROM scratch

LABEL org.opencontainers.image.authors="Mateusz Ł.<example@example.com>"

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/weather /weather

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD ["/weather", "--health"]

EXPOSE 8080

ENTRYPOINT ["/weather"]
