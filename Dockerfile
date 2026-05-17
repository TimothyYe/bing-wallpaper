# syntax=docker/dockerfile:1.6

FROM golang:1.26-alpine AS builder
WORKDIR /src

# Download deps in a separate layer so source-only changes don't bust
# the dependency cache.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO is not needed (pure-Go binary), so produce a static binary that
# can run in a minimal runtime image. -trimpath/-s/-w shrink the binary
# and make builds reproducible. -tags timetzdata embeds the IANA tz
# database so TZ works without zoneinfo files in the runtime image.
RUN CGO_ENABLED=0 go build \
    -trimpath \
    -tags timetzdata \
    -ldflags="-s -w" \
    -o /out/bw ./bw

# Distroless static image: ~5 MB, includes ca-certificates, runs as a
# non-root user, and has no shell or package manager.
FROM gcr.io/distroless/static-debian12:nonroot

ENV TZ=Asia/Shanghai
COPY --from=builder /out/bw /bw

EXPOSE 9000
ENTRYPOINT ["/bw", "run"]
