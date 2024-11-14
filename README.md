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
        root@ubuntu:~/GolandProjects/wallet-service# docker-compose up -d --build 
        [+] Building 7.7s (17/17) FINISHED                                                                                                                                                           docker:default
        => [wallet-service internal] load build definition from Dockerfile                                                                                                                                    0.0s
        => => transferring dockerfile: 866B                                                                                                                                                                   0.0s
        => WARN: FromAsCasing: 'as' and 'FROM' keywords' casing do not match (line 2)                                                                                                                         0.0s
        => [wallet-service internal] load metadata for docker.io/library/alpine:3.19                                                                                                                          0.0s
        => [wallet-service internal] load metadata for docker.io/library/golang:1.23-alpine                                                                                                                   0.0s
        => [wallet-service internal] load .dockerignore                                                                                                                                                       0.0s
        => => transferring context: 2B                                                                                                                                                                        0.0s
        => [wallet-service builder 1/5] FROM docker.io/library/golang:1.23-alpine                                                                                                                             0.0s
        => [wallet-service internal] load build context                                                                                                                                                       0.1s
        => => transferring context: 440.21kB                                                                                                                                                                  0.1s
        => [wallet-service stage-1 1/5] FROM docker.io/library/alpine:3.19                                                                                                                                    0.0s
        => CACHED [wallet-service builder 2/5] WORKDIR /app                                                                                                                                                   0.0s
        => [wallet-service builder 3/5] COPY . .                                                                                                                                                              0.3s
        => [wallet-service builder 4/5] RUN ./bin/golangci-lint run --config .golangci.yml                                                                                                                    6.5s
        => [wallet-service builder 5/5] RUN go build -o wallet-service ./cmd/main.go                                                                                                                          0.8s
        => CACHED [wallet-service stage-1 2/5] WORKDIR /app                                                                                                                                                   0.0s 
        => CACHED [wallet-service stage-1 3/5] COPY --from=builder /app/wallet-service .                                                                                                                      0.0s 
        => CACHED [wallet-service stage-1 4/5] COPY --from=builder /app/config config                                                                                                                         0.0s 
        => CACHED [wallet-service stage-1 5/5] COPY --from=builder /app/pkg pkg                                                                                                                               0.0s 
        => [wallet-service] exporting to image                                                                                                                                                                0.0s 
        => => exporting layers                                                                                                                                                                                0.0s
        => => writing image sha256:ea34bba2bfa75cfa1ad7920b4c9207630a4393345ad743c8d3308ffa53155887                                                                                                           0.0s
        => => naming to docker.io/library/wallet-service-wallet-service                                                                                                                                       0.0s
        => [wallet-service] resolving provenance for metadata file                                                                                                                                            0.0s
        [+] Running 3/0
        ✔ Container wallet-redis     Running                                                                                                                                                                  0.0s 
        ✔ Container wallet-postgres  Running                                                                                                                                                                  0.0s 
        ✔ Container wallet-service   Running                                                                                                                                                                  0.0s 
        root@ubuntu:~/GolandProjects/wallet-service# docker ps 
        CONTAINER ID   IMAGE                           COMMAND                   CREATED         STATUS         PORTS                                       NAMES
        6dec7ac8aa63   wallet-service-wallet-service   "./wallet-service"        6 minutes ago   Up 6 minutes   0.0.0.0:8080->8080/tcp, :::8080->8080/tcp   wallet-service
        abe7f2b74452   postgres:17.0-alpine3.19        "docker-entrypoint.s…"   6 minutes ago   Up 6 minutes   0.0.0.0:5432->5432/tcp, :::5432->5432/tcp   wallet-postgres
        3e7cb405bfe3   redis:alpine                    "docker-entrypoint.s…"   6 minutes ago   Up 6 minutes   0.0.0.0:6379->6379/tcp, :::6379->6379/tcp   wallet-redis
        root@ubuntu:~/GolandProjects/wallet-service# docker-compose ps 
        NAME              IMAGE                           COMMAND                   SERVICE          CREATED         STATUS         PORTS
        wallet-postgres   postgres:17.0-alpine3.19        "docker-entrypoint.s…"   postgres         6 minutes ago   Up 6 minutes   0.0.0.0:5432->5432/tcp, :::5432->5432/tcp
        wallet-redis      redis:alpine                    "docker-entrypoint.s…"   redis            6 minutes ago   Up 6 minutes   0.0.0.0:6379->6379/tcp, :::6379->6379/tcp
        wallet-service    wallet-service-wallet-service   "./wallet-service"        wallet-service   6 minutes ago   Up 6 minutes   0.0.0.0:8080->8080/tcp, :::8080->8080/tcp
    
    // 查看日志
    docker-compose logs -f wallet-service
    ```

3. 使用golang命令行运行运行服务:
    ```bash
    docker-compose stop wallet-service
    go run cmd/main.go

    root@ubuntu:~/GolandProjects/wallet-service#  docker-compose stop wallet-service
    [+] Stopping 1/1
    ✔ Container wallet-service  Stopped                                                                                                                                                                   0.2s 
    root@ubuntu:~/GolandProjects/wallet-service# go run cmd/main.go
    {"level":"info","time":"2024-11-15 00:14:00.230","caller":"/root/GolandProjects/wallet-service/pkg/config/config.go:58","msg":"update config","content":"wallet_service:\n  host: \"127.0.0.1\"\n  port: 8080\n\npostgres:\n  host: \"127.0.0.1\"\n  port: 5432\n  user: \"postgres\"\n  password: \"postgres\"\n  database: \"wallet\"\n  ssl_mode: \"disable\"\n\nredis:\n  host: \"127.0.0.1\"\n  port: 6379\n  password: \"\"\n  db: 0\n  pool_size: 10\n  pool_timeout: 5\n  min_idle_conns: 2\n  max_idle_conns: 5\n  conn_max_idle_time: 300","obj":{"WalletService":{"Port":8080,"Host":"127.0.0.1"},"Postgres":{"Host":"127.0.0.1","Port":5432,"User":"postgres","Password":"postgres","Database":"wallet","SSLMode":"disable"},"Redis":{"Host":"127.0.0.1","Port":6379,"Password":"","DB":0,"PoolSize":10,"PoolTimeout":5,"MinIdleConns":2,"MaxIdleConns":5,"ConnMaxIdleTime":300}}}
    [GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

    [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
    - using env:	export GIN_MODE=release
    - using code:	gin.SetMode(gin.ReleaseMode)

    [GIN-debug] POST   /wallet/:user_id/deposit  --> wallet-service/controllers.(*WalletController).Deposit-fm (3 handlers)
    [GIN-debug] POST   /wallet/:user_id/withdraw --> wallet-service/controllers.(*WalletController).Withdraw-fm (3 handlers)
    [GIN-debug] POST   /wallet/transfer/:sender_id/to/:receiver_id --> wallet-service/controllers.(*WalletController).Transfer-fm (3 handlers)
    [GIN-debug] GET    /wallet/:user_id/balance  --> wallet-service/controllers.(*WalletController).GetBalance-fm (3 handlers)
    [GIN-debug] GET    /wallet/:user_id/transactions --> wallet-service/controllers.(*WalletController).GetTransactionHistory-fm (3 handlers)
    [GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
    Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
    [GIN-debug] Listening and serving HTTP on :8080

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
    go test ./... -race -cover
            wallet-service/cmd              coverage: 0.0% of statements
    ok      wallet-service/controllers      (cached)        coverage: 78.6% of statements
    ?       wallet-service/models   [no test files]
    ok      wallet-service/pkg/config       (cached)        coverage: 84.4% of statements
    ok      wallet-service/pkg/logger       (cached)        coverage: 77.9% of statements
    ok      wallet-service/pkg/postgresx    (cached)        coverage: 83.3% of statements
    ok      wallet-service/pkg/rdsLimit     (cached)        coverage: 69.2% of statements
    ok      wallet-service/pkg/redisx       (cached)        coverage: 77.3% of statements
    ok      wallet-service/services (cached)        coverage: 80.3% of statements
```
## golangci-lint 

```bash
    ./bin/golangci-lint run --config .golangci.yml
    {"Issues":[],"Report":{"Linters":[{"Name":"asasalint"},{"Name":"asciicheck"},{"Name":"bidichk"},{"Name":"bodyclose"},{"Name":"canonicalheader"},{"Name":"containedctx"},{"Name":"contextcheck"},{"Name":"copyloopvar"},{"Name":"cyclop"},{"Name":"decorder"},{"Name":"deadcode"},{"Name":"depguard"},{"Name":"dogsled"},{"Name":"dupl"},{"Name":"dupword"},{"Name":"durationcheck"},{"Name":"errcheck","Enabled":true,"EnabledByDefault":true},{"Name":"errchkjson"},{"Name":"errname"},{"Name":"errorlint"},{"Name":"execinquery"},{"Name":"exhaustive"},{"Name":"exhaustivestruct"},{"Name":"exhaustruct"},{"Name":"exportloopref"},{"Name":"forbidigo"},{"Name":"forcetypeassert"},{"Name":"fatcontext"},{"Name":"funlen"},{"Name":"gci"},{"Name":"ginkgolinter"},{"Name":"gocheckcompilerdirectives"},{"Name":"gochecknoglobals"},{"Name":"gochecknoinits"},{"Name":"gochecksumtype"},{"Name":"gocognit"},{"Name":"goconst"},{"Name":"gocritic","Enabled":true},{"Name":"gocyclo"},{"Name":"godot"},{"Name":"godox"},{"Name":"err113"},{"Name":"gofmt"},{"Name":"gofumpt"},{"Name":"goheader"},{"Name":"goimports"},{"Name":"golint"},{"Name":"mnd"},{"Name":"gomnd"},{"Name":"gomoddirectives"},{"Name":"gomodguard"},{"Name":"goprintffuncname"},{"Name":"gosec"},{"Name":"gosimple","Enabled":true,"EnabledByDefault":true},{"Name":"gosmopolitan"},{"Name":"govet","Enabled":true,"EnabledByDefault":true},{"Name":"grouper"},{"Name":"ifshort"},{"Name":"iface"},{"Name":"importas"},{"Name":"inamedparam"},{"Name":"ineffassign","Enabled":true,"EnabledByDefault":true},{"Name":"interfacebloat"},{"Name":"interfacer"},{"Name":"intrange"},{"Name":"ireturn"},{"Name":"lll"},{"Name":"loggercheck"},{"Name":"maintidx"},{"Name":"makezero"},{"Name":"maligned"},{"Name":"mirror"},{"Name":"misspell"},{"Name":"musttag"},{"Name":"nakedret"},{"Name":"nestif"},{"Name":"nilerr"},{"Name":"nilnil"},{"Name":"nlreturn"},{"Name":"noctx"},{"Name":"nonamedreturns"},{"Name":"nosnakecase"},{"Name":"nosprintfhostport"},{"Name":"paralleltest"},{"Name":"perfsprint"},{"Name":"prealloc"},{"Name":"predeclared"},{"Name":"promlinter"},{"Name":"protogetter"},{"Name":"reassign"},{"Name":"recvcheck"},{"Name":"revive"},{"Name":"rowserrcheck"},{"Name":"sloglint"},{"Name":"scopelint"},{"Name":"sqlclosecheck"},{"Name":"spancheck"},{"Name":"staticcheck","Enabled":true,"EnabledByDefault":true},{"Name":"structcheck"},{"Name":"stylecheck"},{"Name":"tagalign"},{"Name":"tagliatelle"},{"Name":"tenv"},{"Name":"testableexamples"},{"Name":"testifylint"},{"Name":"testpackage"},{"Name":"thelper"},{"Name":"tparallel"},{"Name":"typecheck","Enabled":true,"EnabledByDefault":true},{"Name":"unconvert"},{"Name":"unparam"},{"Name":"unused","Enabled":true,"EnabledByDefault":true},{"Name":"usestdlibvars"},{"Name":"varcheck"},{"Name":"varnamelen"},{"Name":"wastedassign"},{"Name":"whitespace"},{"Name":"wrapcheck"},{"Name":"wsl"},{"Name":"zerologlint"},{"Name":"nolintlint"}]}}
```