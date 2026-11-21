# LoopForge for Railway. The app is upstream's official image, pinned by
# version and digest; this file only adds a start script and a small helper
# that applies upstream's starter SQL (the image has no mysql client).
FROM golang:1.26.11-bookworm@sha256:dc40a13b61304647ed7e5d581d246b86413dfc29d63c6d018310feb93cd98fc0 AS lfdb
WORKDIR /src
COPY lfdb/go.mod lfdb/go.sum ./
RUN go mod download
COPY lfdb/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/lfdb .

FROM ghcr.io/tencent/loopforge:0.4.0@sha256:87b6a3e3cba201d03f8f8d22929881f576021fc78c7c42e82f78329abbf4bda4
ENV LOOPFORGE_VERSION=0.4.0
COPY --from=lfdb /out/lfdb /usr/local/bin/lfdb
# Unchanged copies of upstream v0.4.0 migrations/seed/*.sql, 001 and 031.
COPY sql/ /opt/loopforge-railway/sql/
COPY start.sh /usr/local/bin/loopforge-railway-start
RUN chmod 755 /usr/local/bin/loopforge-railway-start
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/loopforge-railway-start"]
