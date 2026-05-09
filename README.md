# WeBook

基于 Go 的用户系统后端服务，实现了邮箱注册登录、手机号/邮箱验证码登录、用户信息管理等核心功能。

## 技术栈

| 类别 | 技术 |
|------|------|
| 语言 | Go |
| Web 框架 | Gin |
| ORM | GORM |
| 数据库 | MySQL |
| 缓存 | Redis |
| 认证 | JWT（HS512） |
| 短信服务 | 腾讯云短信 SDK |
| 邮件服务 | QQ 邮箱 SMTP |
| 容器化 | Docker / docker-compose |
| 部署 | Kubernetes + Nginx Ingress |

## 项目结构

```
webook/
├── cmd/
│   └── server/          # 程序入口
├── config/              # 环境配置（dev / k8s）
├── internal/
│   ├── domain/          # 业务实体定义
│   ├── handler/         # HTTP 处理层，参数校验与响应
│   ├── middleware/       # Gin 中间件（JWT 校验、CORS、登录态）
│   ├── repository/      # 数据访问层
│   │   ├── cache/       # Redis 缓存（用户信息、验证码 + Lua 脚本）
│   │   └── dao/         # 数据库操作
│   ├── router/          # 路由注册
│   └── service/         # 业务逻辑层
│       └── sms/         # 短信/邮件抽象接口及实现
├── docker-compose.yaml  # 本地开发环境
├── Dockerfile
└── k8s-*.yaml           # Kubernetes 部署配置
```

## 核心功能

- **邮箱注册**：正则校验邮箱格式与密码强度，bcrypt 加盐加密存储密码
- **邮箱登录**：JWT 签发与校验，Token 不足 7 天自动续签，比对 User-Agent 防盗用
- **验证码登录**：支持手机号短信与邮箱两种渠道，Redis + Lua 脚本原子控制发送频率与校验次数
- **用户信息**：查看与编辑个人资料（昵称、生日、简介），查询时走 Redis 缓存，未命中再查 MySQL
- **登出**：清除 Session

## 本地运行

**前置依赖**：Go 1.21+、Docker

```bash
# 1. 启动 MySQL 和 Redis
docker-compose up -d

# 2. 运行服务
go run cmd/server/mian.go
```

服务默认监听 `:8082`。

## 接口列表

| 方法 | 路径 | 说明 | 是否需要登录 |
|------|------|------|------------|
| POST | /users/signup | 邮箱注册 | 否 |
| POST | /users/login | 邮箱登录（JWT） | 否 |
| POST | /users/logout | 登出 | 是 |
| POST | /users/edit | 编辑个人信息 | 是 |
| GET | /users/profile | 查看个人信息 | 是 |

> 需要登录的接口请在请求头携带 `Authorization: Bearer <token>`，登录成功后 token 从响应头 `x-jwt-token` 获取。

## 部署（Kubernetes）

```bash
# 部署 MySQL、Redis 及应用服务
kubectl apply -f k8s-mysql-pv.yaml
kubectl apply -f k8s-mysql-pvc.yaml
kubectl apply -f k8s-mysql-deployment.yaml
kubectl apply -f k8s-mysql-service.yaml
kubectl apply -f k8s-redis-deployment.yaml
kubectl apply -f k8s-redis-service.yaml
kubectl apply -f k8s-webook-deployment.yaml
kubectl apply -f k8s-webook-service.yaml
kubectl apply -f k8s-ingress-niginx.yaml
```

应用默认部署 3 个副本，通过 Nginx Ingress 做负载均衡。

## 设计说明

- **分层架构**：handler → service → repository → dao，各层只依赖下层接口，便于替换实现
- **短信抽象**：`sms.Service` 接口统一屏蔽腾讯云短信与邮箱两种实现，业务层无感知切换
- **缓存策略**：用户信息查询优先读 Redis（TTL 1小时），缓存未命中回源 MySQL 并写入缓存
- **原子操作**：验证码的发送频率限制与错误次数锁定通过 Lua 脚本在 Redis 侧原子执行，避免并发竞态

## License

MIT