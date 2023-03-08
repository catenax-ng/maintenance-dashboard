FROM golang:1.20-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go test ./... && \
    CGO_ENABLED=0 \
    go build -installsuffix 'static' -ldflags="-w -s" .


FROM gcr.io/distroless/static:nonroot AS final

WORKDIR /app

COPY --from=builder --chown=nonroot:nonroot /app/maintenance-dashboard /app/maintenance-dashboard

HEALTHCHECK --interval=5m --timeout=3s \
  CMD curl -f http://localhost:2112/health || exit 1

ENTRYPOINT ["/app/maintenance-dashboard"]

CMD ["-in-cluster=true"]