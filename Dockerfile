# =======================================================================================
ARG GO_VERSION="1.25.1-alpine3.22"
ARG DISTROLESS_VERSION="nonroot"
# =======================================================================================
FROM golang:${GO_VERSION} AS builder
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/* \
    && rm -rf /tmp/* \
    && rm -rf /var/tmp/*
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -a -installsuffix cgo \
    -ldflags="-w -s -X main.version=$(date +%Y%m%d-%H%M%S) -extldflags '-static'" \
    -trimpath \
    -buildvcs=false \
    -o /app/bin .
# =======================================================================================
FROM gcr.io/distroless/static-debian12:${DISTROLESS_VERSION}
ENV PORT=8080
ENV GOGC=100
ENV GOMEMLIMIT=128MiB
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY --from=builder /app/bin /app/bin
COPY --from=builder /app/migrations /app/migrations
USER nonroot:nonroot
EXPOSE ${PORT}
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 CMD ["/app/bin", "healthcheck"] || exit 1
ENTRYPOINT ["/app/bin"]
CMD ["serve"]
# =======================================================================================
