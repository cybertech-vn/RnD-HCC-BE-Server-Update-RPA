# Stage 1: Build environment
FROM golang:alpine AS builder

# Cài dependencies cho CGO
# Thêm tzdata để copy sang runtime
RUN apk add --no-cache gcc musl-dev tzdata

WORKDIR /app

# Tối ưu cache: Copy go.mod trước để tận dụng Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ source code
COPY . .

# Build server + passkey (Tối ưu dung lượng với -ldflags="-s -w")
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o setup-server ./cmd/server

# Stage 2: Runtime
FROM alpine:latest

# Tạo non-root user để tăng tính bảo mật
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Cài chứng chỉ SSL và Timezone
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binaries và timezone từ builder
COPY --from=builder /app/setup-server .
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# KHÔNG copy config.yml ở đây vì docker-compose đã mount volume file này

# Tạo thư mục log và phân quyền cho appuser
RUN mkdir -p ./log && chown -R appuser:appgroup /app

# Chuyển sang dùng non-root user
USER appuser

EXPOSE 7878

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:7878/heath || exit 1

CMD ["./setup-server"]
