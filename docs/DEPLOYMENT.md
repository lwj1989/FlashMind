# 部署指南

本文档介绍如何将 FlashMind 部署到不同的环境中。

## 🚀 部署选项

### 1. 本地部署

#### 使用脚本部署（推荐）

```bash
# 克隆项目
git clone https://github.com/lwj1989/FlashMind.git
cd FlashMind

# 一键启动
./start.sh
```

访问应用：
- 前端：http://localhost:5173
- 后端：http://localhost:8080

#### 手动部署

```bash
# 后端部署
cd backend
go mod download
go build -o flashcard cmd/server/main.go
./flashcard

# 前端部署
cd frontend
npm install
npm run build
npm run preview
```

### 2. Docker 部署

#### 使用 Docker Compose（推荐）

```bash
# 克隆项目
git clone https://github.com/lwj1989/FlashMind.git
cd FlashMind

# 启动服务
docker-compose up -d

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

#### 使用单独的 Docker 容器

```bash
# 构建镜像
docker build -t flashmind .

# 运行容器
docker run -d \
  --name flashmind \
  -p 8080:8080 \
  -v flashmind_data:/app/data \
  flashmind

# 查看日志
docker logs -f flashmind
```

### 3. 云平台部署

#### AWS 部署

##### 使用 EC2

```bash
# 1. 创建 EC2 实例（Ubuntu 20.04 LTS）
# 2. 安装 Docker
sudo apt update
sudo apt install docker.io docker-compose -y
sudo systemctl start docker
sudo systemctl enable docker

# 3. 部署应用
git clone https://github.com/lwj1989/FlashMind.git
cd FlashMind
sudo docker-compose up -d

# 4. 配置安全组
# 开放端口 80, 443, 8080
```

##### 使用 ECS

```yaml
# ecs-task-definition.json
{
  "family": "flashmind",
  "taskRoleArn": "arn:aws:iam::ACCOUNT:role/ecsTaskRole",
  "executionRoleArn": "arn:aws:iam::ACCOUNT:role/ecsTaskExecutionRole",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "containerDefinitions": [
    {
      "name": "flashmind",
      "image": "your-registry/flashmind:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/flashmind",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

#### Google Cloud Platform 部署

##### 使用 Cloud Run

```bash
# 1. 构建并推送镜像
gcloud builds submit --tag gcr.io/PROJECT-ID/flashmind

# 2. 部署到 Cloud Run
gcloud run deploy flashmind \
  --image gcr.io/PROJECT-ID/flashmind \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

##### 使用 GKE

```yaml
# kubernetes/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: flashmind
spec:
  replicas: 3
  selector:
    matchLabels:
      app: flashmind
  template:
    metadata:
      labels:
        app: flashmind
    spec:
      containers:
      - name: flashmind
        image: gcr.io/PROJECT-ID/flashmind:latest
        ports:
        - containerPort: 8080
        env:
        - name: GIN_MODE
          value: "release"
---
apiVersion: v1
kind: Service
metadata:
  name: flashmind-service
spec:
  selector:
    app: flashmind
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

#### Azure 部署

##### 使用 Container Instances

```bash
# 创建资源组
az group create --name flashmind-rg --location eastus

# 部署容器
az container create \
  --resource-group flashmind-rg \
  --name flashmind \
  --image your-registry/flashmind:latest \
  --dns-name-label flashmind-app \
  --ports 8080
```

### 4. VPS 部署

#### 使用 Nginx 反向代理

```nginx
# /etc/nginx/sites-available/flashmind
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# 启用站点
sudo ln -s /etc/nginx/sites-available/flashmind /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

#### 使用 SSL/TLS

```bash
# 安装 Certbot
sudo apt install certbot python3-certbot-nginx

# 获取 SSL 证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo crontab -e
# 添加：0 12 * * * /usr/bin/certbot renew --quiet
```

## 🔧 环境配置

### 环境变量

创建 `.env` 文件：

```bash
# 应用配置
GIN_MODE=release
PORT=8080

# 数据库配置
SQLITE_DB_PATH=/app/data/flashcard.db

# 静态文件配置
STATIC_PATH=/app/static

# 安全配置
JWT_SECRET=your-jwt-secret
CORS_ORIGIN=https://your-domain.com

# 日志配置
LOG_LEVEL=info
LOG_FILE=/app/logs/app.log
```

### 生产环境优化

#### 性能优化

```bash
# 后端优化
export GOMAXPROCS=4
export GOGC=100

# 数据库优化
# SQLite 配置
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA cache_size = 1000000;
PRAGMA foreign_keys = true;
PRAGMA temp_store = memory;
```

#### 安全配置

```bash
# 防火墙配置
sudo ufw allow ssh
sudo ufw allow 80
sudo ufw allow 443
sudo ufw enable

# 限制文件权限
chmod 600 .env
chmod 755 flashcard
```

## 📊 监控和日志

### 日志配置

```bash
# 使用 systemd 管理服务
# /etc/systemd/system/flashmind.service
[Unit]
Description=FlashMind Application
After=network.target

[Service]
Type=simple
User=flashmind
WorkingDirectory=/opt/flashmind
ExecStart=/opt/flashmind/flashcard
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target

# 启用服务
sudo systemctl enable flashmind
sudo systemctl start flashmind
```

### 监控设置

#### 使用 Prometheus + Grafana

```yaml
# docker-compose.monitoring.yml
version: '3.8'
services:
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
```

#### 健康检查

```go
// 在后端添加健康检查端点
func HealthCheck(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "ok",
        "timestamp": time.Now(),
        "version": "1.0.0",
    })
}
```

## 🔄 备份和恢复

### 数据库备份

```bash
# 自动备份脚本
#!/bin/bash
# backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/opt/flashmind/backups"
DB_PATH="/opt/flashmind/data/flashcard.db"

mkdir -p $BACKUP_DIR

# 备份数据库
cp $DB_PATH $BACKUP_DIR/flashcard_$DATE.db

# 压缩旧备份
gzip $BACKUP_DIR/flashcard_$DATE.db

# 删除7天前的备份
find $BACKUP_DIR -name "*.gz" -mtime +7 -delete

echo "Backup completed: flashcard_$DATE.db.gz"
```

```bash
# 设置定时备份
crontab -e
# 添加：0 2 * * * /opt/flashmind/backup.sh
```

### 数据恢复

```bash
# 恢复数据库
#!/bin/bash
# restore.sh

BACKUP_FILE=$1
DB_PATH="/opt/flashmind/data/flashcard.db"

if [ -z "$BACKUP_FILE" ]; then
    echo "Usage: $0 <backup_file>"
    exit 1
fi

# 停止服务
sudo systemctl stop flashmind

# 备份当前数据库
cp $DB_PATH ${DB_PATH}.backup

# 恢复数据库
gunzip -c $BACKUP_FILE > $DB_PATH

# 启动服务
sudo systemctl start flashmind

echo "Database restored from $BACKUP_FILE"
```

## 🚦 负载均衡和高可用

### 使用 Nginx 负载均衡

```nginx
# upstream 配置
upstream flashmind_backend {
    server 127.0.0.1:8080;
    server 127.0.0.1:8081;
    server 127.0.0.1:8082;
}

server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://flashmind_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

### Docker Swarm 部署

```yaml
# docker-stack.yml
version: '3.8'
services:
  flashmind:
    image: flashmind:latest
    deploy:
      replicas: 3
      update_config:
        parallelism: 1
        delay: 10s
      restart_policy:
        condition: on-failure
    ports:
      - "8080:8080"
    networks:
      - flashmind_network

networks:
  flashmind_network:
    driver: overlay
```

```bash
# 部署到 Swarm
docker stack deploy -c docker-stack.yml flashmind
```

## 🔍 故障排除

### 常见问题

#### 服务无法启动
```bash
# 检查日志
journalctl -u flashmind -f

# 检查端口占用
sudo netstat -tlnp | grep :8080

# 检查文件权限
ls -la /opt/flashmind/
```

#### 数据库连接问题
```bash
# 检查数据库文件
ls -la /opt/flashmind/data/

# 检查数据库权限
sudo -u flashmind sqlite3 /opt/flashmind/data/flashcard.db ".tables"
```

#### 性能问题
```bash
# 监控系统资源
htop
iotop
df -h

# 检查应用性能
curl -w "@curl-format.txt" -o /dev/null -s "http://localhost:8080/api/v1/system/health"
```

### 性能调优

```bash
# 系统级优化
echo 'vm.swappiness=10' >> /etc/sysctl.conf
echo 'net.core.somaxconn=65535' >> /etc/sysctl.conf
sysctl -p

# 文件描述符限制
echo 'flashmind soft nofile 65535' >> /etc/security/limits.conf
echo 'flashmind hard nofile 65535' >> /etc/security/limits.conf
```

## 📞 获取支持

如果在部署过程中遇到问题：

1. 查看 [故障排除指南](../README.md#故障排除)
2. 创建 [GitHub Issue](https://github.com/lwj1989/FlashMind/issues)
3. 查看 [讨论区](https://github.com/lwj1989/FlashMind/discussions)

## 📋 部署检查清单

### 部署前检查
- [ ] 确认系统要求
- [ ] 准备环境变量
- [ ] 配置域名和SSL
- [ ] 设置防火墙规则
- [ ] 准备监控和日志

### 部署后验证
- [ ] 应用正常启动
- [ ] API 端点响应正常
- [ ] 前端页面正常加载
- [ ] 数据库连接正常
- [ ] SSL 证书有效
- [ ] 监控和日志正常

### 维护任务
- [ ] 设置自动备份
- [ ] 配置监控告警
- [ ] 设置自动更新
- [ ] 定期安全检查
