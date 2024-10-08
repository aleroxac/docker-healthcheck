FROM golang:1.23.0 AS base

ARG \
    HEALTHCHECK_PROTOCOL \
    HEALTHCHECK_HOST \
    HEALTHCHECK_PORT \
    HEALTHCHECK_PATH \
    API_PORT

ENV API_PORT=${API_PORT}

FROM base AS build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o server main.go
RUN go install "github.com/aleroxac/docker-healthcheck/cmd/healthcheck/main.go"
RUN pwd; which healthcheck

# FROM scratch
FROM base AS main
WORKDIR /app
COPY --from=build /build/server /go/bin/healthcheck /app/
CMD ["/app/server"]
EXPOSE ${API_PORT}

ENV \
    HEALTHCHECK_PROTOCOL=${HEALTHCHECK_PROTOCOL} \
    HEALTHCHECK_HOST=${HEALTHCHECK_HOST} \
    HEALTHCHECK_PORT=${HEALTHCHECK_PORT} \
    HEALTHCHECK_PATH=${HEALTHCHECK_PATH}
HEALTHCHECK \
    --start-period=2s \
    --interval=10s \
    --timeout=5s \
    --retries=3 \
    CMD [ "/app/healthcheck" ]

