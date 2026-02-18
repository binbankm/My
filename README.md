# ServerPanel - Linux服务器管理面板

一个轻量级的Linux服务器管理面板，类似1Panel，提供系统监控、容器管理、文件管理等功能。

## 特性

- 🚀 单二进制部署，开箱即用
- 📊 实时系统监控（CPU、内存、磁盘、网络）
- 🐳 Docker容器管理
- 📁 文件管理器
- 🗄️ 数据库管理
- 🌐 Nginx/Web服务器管理
- 🔐 安全认证和用户管理
- 💻 前端: Vue 3 + shadcn/ui
- ⚡ 后端: Go + Gin

## 快速开始

### 安装

```bash
# 下载最新版本
wget https://github.com/binbankm/My/releases/latest/download/serverpanel-linux-amd64.tar.gz

# 解压
tar -zxvf serverpanel-linux-amd64.tar.gz

# 安装
cd serverpanel
sudo ./install.sh

# 启动服务
sudo systemctl start serverpanel
```

### 访问

打开浏览器访问: `http://your-server-ip:8888`

默认账号: `admin`  
默认密码: `admin123`

## 开发

### 前置要求

- Go 1.21+
- Node.js 18+
- npm/pnpm

### 开发模式

```bash
# 启动后端
go run cmd/server/main.go

# 启动前端（新终端）
cd frontend
npm install
npm run dev
```

### 构建

```bash
# 构建所有平台
make build-all

# 仅构建Linux
make build-linux
```

## 文档

详见 [docs](./docs) 目录

## 许可证

MIT License
