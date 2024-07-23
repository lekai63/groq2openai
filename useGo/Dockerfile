# 构建阶段
FROM hub.atomgit.com/library/golang:1-bullseye AS builder

# 设置工作目录
WORKDIR /app

# 复制源代码
COPY main.go .

ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on

# 初始化 Go 模块
RUN go mod init go-audio-transcription-server
RUN go mod tidy

# 下载依赖
RUN go mod download

# 构建 Linux amd64 二进制文件
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -a -installsuffix cgo -o main .

# 运行阶段
FROM hub.atomgit.com/library/debian:bullseye-slim

# 安装 CA 证书
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8000

# 运行
CMD ["./main"]
