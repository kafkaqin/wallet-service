# Wallet Service

该钱包服务基于Go和Gin框架实现，使用PostgreSQL存储用户和交易数据，Redis作为缓存数据库。服务提供了RESTful API，支持以下功能：
- 存款
- 取款
- 转账
- 查询余额
- 查询交易历史

## 技术栈
- 语言: Go 1.23.1
- 框架: Gin
- 数据库: PostgreSQL
- 缓存数据库: Redis

## 快速开始

### 本地运行

1. 克隆项目:
    ```bash
    git clone https://github.com/kafkaqin/wallet-service.git
    cd wallet-service
    ```

2. 在安装docker和docker-compose,启动数据库,钱包服务:
    - 使用Docker Compose启动PostgreSQL和Redis。
    - 使用docker-compose运行
    ```bash

    docker-compose up --build -d

    root@ubuntu:~/GolandProjects/wallet-service# docker-compose ps 
    WARN[0000] /root/GolandProjects/wallet-service/docker-compose.yaml: the attribute `version` is obsolete, it will be ignored, please remove it to avoid potential confusion 
    NAME              IMAGE                           COMMAND                   SERVICE          CREATED          STATUS          PORTS
    wallet-postgres   postgres:17.0-alpine3.19        "docker-entrypoint.s…"   postgres         25 hours ago     Up 15 minutes   0.0.0.0:5432->5432/tcp, :::5432->5432/tcp
    wallet-redis      redis:alpine                    "docker-entrypoint.s…"   redis            25 hours ago     Up 15 minutes   0.0.0.0:6379->6379/tcp, :::6379->6379/tcp
    wallet-service    wallet-service-wallet-service   "./wallet-service"        wallet-service   11 minutes ago   Up 11 minutes   0.0.0.0:8080->8080/tcp, :::8080->8080/tcp
    
    // 查看日志
    docker-compose logs -f wallet-service
    ```

3. 使用golang命令行运行运行服务:
    ```bash
    docker-compose stop wallet-service
    go run cmd/main.go
    ```

### API 端点

- `POST /wallet/:user_id/deposit`: 向指定用户钱包存入金额。
- `POST /wallet/:user_id/withdraw`: 从指定用户钱包取出金额。
- `POST /wallet/transfer/:sender_id/to/:receiver_id`: 从一个用户钱包转账到另一个用户钱包。
- `GET /wallet/:user_id/balance`: 查询指定用户钱包的余额。
- `GET /wallet/:user_id/transactions`: 查询指定用户的交易历史。

- postman文件 postman/wallet-service.postman_collection.json
```
root@ubuntu:~/GolandProjects/wallet-service# ls -l postman/
total 4
-rw-r--r-- 1 root root 2188 11月 13 19:36 wallet-service.postman_collection.json
```
## 测试
使用以下命令运行测试并生成覆盖率报告:
```bash
    go test ./... -race -coverprofile=coverage.out
            wallet-service/cmd              coverage: 0.0% of statements
    ?       wallet-service/models   [no test files]
    ok      wallet-service/controllers      (cached)        coverage: 62.4% of statements
    ok      wallet-service/pkg/config       (cached)        coverage: 66.7% of statements
    ok      wallet-service/pkg/logger       (cached)        coverage: 78.4% of statements
    ok      wallet-service/pkg/postgresx    (cached)        coverage: 84.9% of statements
    ok      wallet-service/pkg/rdsLimit     (cached)        coverage: 69.2% of statements
    ok      wallet-service/pkg/redisx       (cached)        coverage: 68.2% of statements
```
