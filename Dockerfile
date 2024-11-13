# 使用golang镜像作为基础镜像
FROM golang:1.23-alpine as builder

# 设置工作目录
WORKDIR /app

# 安装golangci-lint
# RUN #curl -sSf https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s

# 复制项目文件到容器
COPY . .

# 运行golangci-lint
# RUN ./bin/golangci-lint run --config .golangci.yml

# 编译Go应用
RUN go build -o wallet-service ./cmd/main.go

# 第二阶段：创建最终镜像
FROM alpine:3.19

# 安装必要的库
#RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /app

# 将构建的二进制文件从构建阶段复制到最终镜像
COPY --from=builder /app/wallet-service .
COPY --from=builder /app/config config
COPY --from=builder /app/pkg pkg
# 暴露端口
EXPOSE 8080

# 设置默认命令启动应用
CMD ["./wallet-service"]
