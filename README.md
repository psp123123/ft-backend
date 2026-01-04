# ft-backend

一个高性能的服务器管理平台后端项目，基于 Go + Gin + GORM + MySQL 构建，提供机器管理、Kubernetes 部署、API 服务等功能。

## 技术栈

- **语言**: Go 1.24.5
- **Web 框架**: Gin 1.11.0
- **ORM**: GORM 1.31.1
- **数据库**: MySQL
- **配置管理**: YAML
- **认证**: JWT
- **实时通信**: WebSocket
- **加密**: Golang Crypto

## 功能特性

### 核心功能
- **API 服务**: 提供完整的 RESTful API 接口
- **认证授权**: JWT 认证和基于角色的权限控制
- **机器管理**: 服务器信息管理和资源监控
- **Kubernetes 部署**: 集群管理和版本控制
- **文件管理**: 文件上传、下载和共享功能

### 高级功能
- **安全审计**: 操作日志记录和审计
- **实时通信**: WebSocket 服务支持
- **数据传输**: 大文件传输支持
- **跨域支持**: CORS 中间件配置
- **数据库迁移**: 自动数据库模式迁移

## 快速开始

### 环境要求
- Go >= 1.24.0
- MySQL >= 5.7.0

### 安装依赖
```bash
go mod download
```

### 配置数据库
1. 创建 MySQL 数据库
2. 修改 `config.yaml` 文件中的数据库配置：

```yaml
database:
  host: localhost
  port: 3306
  user: root
  password: your_password
  dbname: ft_backend
```

### 初始化数据库
执行 SQL 初始化脚本：
```bash
mysql -u root -p ft_backend < database/schema.sql
mysql -u root -p ft_backend < database/init_k8s_versions.sql
```

### 编译和运行
```bash
# 编译
go build -o ft-backend

# 运行
./ft-backend
```

或者直接运行：
```bash
go run main.go
```

服务将在 `http://localhost:8000` 启动。

## 项目结构

```
ft-backend/
├── config/                # 配置文件
├── database/              # 数据库相关
│   ├── database.go        # 数据库连接配置
│   ├── schema.sql         # 数据库初始化脚本
│   └── init_k8s_versions.sql # Kubernetes 版本初始化
├── handlers/              # API 处理器
│   ├── auth.go            # 认证相关
│   ├── dashboard.go       # 仪表盘
│   ├── k8s_deploy.go      # Kubernetes 部署
│   ├── machine.go         # 机器管理
│   ├── user.go            # 用户管理
│   └── websocket.go       # WebSocket 服务
|---iotservice/
|   |--- heatbeat.go        # 客户端client相关
|
├── middleware/            # 中间件
│   ├── auth.go            # 认证中间件
│   └── cors.go            # CORS 中间件
├── models/                # 数据模型
│   ├── k8s_cluster.go     # Kubernetes 集群
│   ├── k8s_version.go     # Kubernetes 版本
│   ├── machine.go         # 机器
│   ├── user.go            # 用户
│   └── operation_log.go   # 操作日志
├── routes/                # 路由配置
│   └── router.go          # 路由注册
├── utils/                 # 工具函数
│   ├── jwt.go             # JWT 工具
│   ├── password.go        # 密码工具
│   └── websocket_manager.go # WebSocket 管理器
├── uploads/               # 上传文件存储
├── config.yaml            # 主配置文件
├── go.mod                 # Go 模块
├── go.sum                 # 依赖校验
├── main.go                # 入口文件
└── README.md              # 项目说明
```

## 配置说明

### 主配置文件 (config.yaml)
```yaml
# 服务器配置
server:
  port: 8000
  mode: debug  # debug/release

# 数据库配置
database:
  host: localhost
  port: 3306
  user: root
  password: your_password
  dbname: ft_backend
  charset: utf8mb4
  max_idle_conns: 10
  max_open_conns: 100

# JWT 配置
jwt:
  secret: your_jwt_secret
  expires_hours: 24

# 上传配置
upload:
  path: ./uploads
  max_size: 104857600  # 100MB
```

## API 文档

### 认证接口
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/refresh` - 刷新令牌

### 机器管理接口
- `GET /api/machines` - 获取机器列表
- `POST /api/machines` - 创建机器
- `GET /api/machines/:id` - 获取机器详情
- `PUT /api/machines/:id` - 更新机器
- `DELETE /api/machines/:id` - 删除机器

### Kubernetes 接口
- `GET /api/k8s/versions` - 获取 Kubernetes 版本列表
- `POST /api/k8s/clusters` - 创建集群
- `GET /api/k8s/clusters` - 获取集群列表

### 用户管理接口
- `GET /api/users` - 获取用户列表
- `POST /api/users` - 创建用户
- `PUT /api/users/:id` - 更新用户
- `DELETE /api/users/:id` - 删除用户

### client客户端接口
- `POST /v1/heartbeats` - 上传心跳信息

发送心跳数据结构
```bash
{
    "client_id": "123qweasd",
    "heartbeat_time": 1735689600000,
    "client_version": "1.0.3",
    "process_id": 12345,
    "status": "normal",
    "local_ip": "192.168.1.5",
    "business_module": "log-collect",
    "task_count": 8,
    "last_task_time": 1735689000000,
    "primary_host": {
        "ip": "192.168.56.11",
        "hostname": "node-1",
        "os_info": "Linux x86_64",
        "cpu_usage": 25.6,
        "memory_usage": 104857600,
        "disk_usage": "200GB",
        "network_delay": 15,
        "network_interface": "eth0",
        "status": "up"
    },
    "secondary_hosts": [
        {
            "ip": "192.168.56.101",
            "hostname": "node-2",
            "os_info": "Linux x86_64",
            "cpu_usage": 25.6,
            "memory_usage": 104857600,
            "disk_usage": "200GB",
            "network_delay": 15,
            "network_interface": "eth0",
            "status": "up"
        },
        {
            "ip": "192.168.56.102",
            "hostname": "node-3",
            "os_info": "Linux x86_64",
            "cpu_usage": 25.6,
            "memory_usage": 104857600,
            "disk_usage": "200GB",
            "network_delay": 15,
            "network_interface": "eth0",
            "status": "up"
        },
        {
            "ip": "192.168.56.103",
            "hostname": "node-4",
            "os_info": "Linux x86_64",
            "cpu_usage": 25.6,
            "memory_usage": 104857600,
            "disk_usage": "200GB",
            "network_delay": 15,
            "network_interface": "eth0",
            "status": "up"
        }
    ]
    }
```

## 开发规范

### 代码规范
- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 使用 `go vet` 检查代码
- 函数名使用驼峰命名法
- 变量名使用驼峰命名法

### 提交规范
- 使用 Conventional Commits 规范
- 提交信息格式：`type(scope): description`

## 开发流程

### 启动开发服务器
```bash
go run main.go
```

### 构建生产版本
```bash
go build -o ft-backend -ldflags "-s -w"
```

### 代码检查
```bash
go fmt ./...
go vet ./...
```

## 数据库管理

### 生成数据库迁移
```bash
# 使用 GORM 自动迁移
# 在 main.go 中已集成自动迁移功能
```

### 更新管理员密码
执行 SQL 脚本：
```bash
mysql -u root -p ft_backend < update_admin_password.sql
```

## 安全注意事项
- 定期更换 JWT 密钥
- 使用 HTTPS 协议部署生产环境
- 限制数据库用户权限
- 定期备份数据库

## 许可证

MIT License
