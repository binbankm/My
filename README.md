# ServerPanel - Linux服务器管理面板

<div align="center">

一个轻量级的Linux服务器管理面板，类似1Panel，提供系统监控、容器管理、文件管理等功能。

A lightweight Linux server management panel similar to 1Panel, providing system monitoring, container management, file management, and more.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![Vue Version](https://img.shields.io/badge/Vue-3.4+-4FC08D?logo=vue.js)](https://vuejs.org/)

[English](#english) | [中文](#chinese)

</div>

---

## <a name="chinese"></a>🇨🇳 中文

### ✨ 特性

- 🚀 **单二进制部署** - 无需复杂配置，一个文件即可运行
- 📊 **实时系统监控** - CPU、内存、磁盘、网络实时监控
- 🐳 **Docker容器管理** - 管理Docker容器和镜像（预留接口）
- 📁 **文件管理器** - 完整的文件管理功能，支持上传下载
- 🗄️ **数据库管理** - 数据库管理界面（支持扩展）
- 🌐 **Web服务器管理** - Nginx/Apache配置管理（预留功能）
- 🔐 **安全认证** - JWT认证 + bcrypt密码加密
- 💻 **现代化UI** - Vue 3 + Tailwind CSS + shadcn/ui 组件
- ⚡ **高性能** - Go后端 + Gin框架
- 🔧 **易于扩展** - 模块化设计，易于添加新功能

### 📦 技术栈

**后端**
- Go 1.21+
- Gin Web Framework
- GORM (SQLite)
- JWT认证
- gopsutil (系统监控)

**前端**
- Vue 3 (Composition API)
- Vue Router
- Pinia (状态管理)
- Tailwind CSS
- shadcn/ui 组件
- Axios

### 🚀 快速开始

#### 安装

**方法一：一键远程安装（推荐）**

```bash
curl -fsSL https://raw.githubusercontent.com/binbankm/My/main/scripts/remote-install.sh | sudo bash
```

或者使用 wget：

```bash
wget -qO- https://raw.githubusercontent.com/binbankm/My/main/scripts/remote-install.sh | sudo bash
```

**方法二：手动安装**

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

#### 访问

打开浏览器访问: `http://your-server-ip:8888`

**默认账号**
- 用户名: `admin`  
- 密码: `admin123`

> ⚠️ **重要**: 首次登录后请立即修改默认密码！

#### 卸载

```bash
# 下载并运行卸载脚本
wget -qO- https://raw.githubusercontent.com/binbankm/My/main/scripts/uninstall.sh | sudo bash
```

或者如果您已经有安装包：

```bash
cd serverpanel
sudo ./uninstall.sh
```

### 📖 文档

- [部署指南](docs/DEPLOYMENT.md) - 详细的安装和配置说明
- [开发指南](docs/DEVELOPMENT.md) - 开发环境设置和贡献指南
- [API文档](docs/API.md) - REST API接口文档

### 🛠️ 开发

```bash
# 克隆仓库
git clone https://github.com/binbankm/My.git
cd My

# 启动后端
go run main.go

# 启动前端（新终端）
cd frontend
npm install
npm run dev
```

### 📦 构建

```bash
# 构建所有平台
make build-all

# 仅构建Linux
make build-linux

# 仅构建前端
cd frontend && npm run build

# 仅构建后端
go build -o serverpanel main.go
```

### 🎯 路线图

- [x] 基础系统监控
- [x] 文件管理
- [x] 用户认证
- [ ] Docker完整集成
- [ ] 数据库连接管理（MySQL/PostgreSQL）
- [ ] Nginx配置管理
- [ ] 定时任务管理
- [ ] 日志查看器
- [ ] 终端/SSH集成
- [ ] 备份和恢复
- [ ] 多用户权限管理
- [ ] WebSocket实时通信

### 🤝 贡献

欢迎贡献！请查看 [DEVELOPMENT.md](docs/DEVELOPMENT.md) 了解如何开始。

### 📝 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

---

## <a name="english"></a>🇬🇧 English

### ✨ Features

- 🚀 **Single Binary Deployment** - No complex configuration, run with one file
- 📊 **Real-time System Monitoring** - CPU, memory, disk, network monitoring
- 🐳 **Docker Container Management** - Manage Docker containers and images (API ready)
- 📁 **File Manager** - Complete file management with upload/download support
- 🗄️ **Database Management** - Database management interface (extensible)
- 🌐 **Web Server Management** - Nginx/Apache configuration (planned)
- 🔐 **Secure Authentication** - JWT auth + bcrypt password hashing
- 💻 **Modern UI** - Vue 3 + Tailwind CSS + shadcn/ui components
- ⚡ **High Performance** - Go backend + Gin framework
- 🔧 **Easy to Extend** - Modular design for easy feature additions

### 📦 Tech Stack

**Backend**
- Go 1.21+
- Gin Web Framework
- GORM (SQLite)
- JWT Authentication
- gopsutil (System monitoring)

**Frontend**
- Vue 3 (Composition API)
- Vue Router
- Pinia (State Management)
- Tailwind CSS
- shadcn/ui Components
- Axios

### 🚀 Quick Start

#### Installation

**Method 1: One-line Remote Installation (Recommended)**

```bash
curl -fsSL https://raw.githubusercontent.com/binbankm/My/main/scripts/remote-install.sh | sudo bash
```

Or using wget:

```bash
wget -qO- https://raw.githubusercontent.com/binbankm/My/main/scripts/remote-install.sh | sudo bash
```

**Method 2: Manual Installation**

```bash
# Download latest release
wget https://github.com/binbankm/My/releases/latest/download/serverpanel-linux-amd64.tar.gz

# Extract
tar -zxvf serverpanel-linux-amd64.tar.gz

# Install
cd serverpanel
sudo ./install.sh

# Start service
sudo systemctl start serverpanel
```

#### Access

Open browser and visit: `http://your-server-ip:8888`

**Default Credentials**
- Username: `admin`  
- Password: `admin123`

> ⚠️ **Important**: Change the default password immediately after first login!

#### Uninstallation

```bash
# Download and run uninstall script
wget -qO- https://raw.githubusercontent.com/binbankm/My/main/scripts/uninstall.sh | sudo bash
```

Or if you already have the package:

```bash
cd serverpanel
sudo ./uninstall.sh
```

### 📖 Documentation

- [Deployment Guide](docs/DEPLOYMENT.md) - Detailed installation and configuration
- [Development Guide](docs/DEVELOPMENT.md) - Development setup and contribution guidelines
- [API Documentation](docs/API.md) - REST API reference

### 🛠️ Development

```bash
# Clone repository
git clone https://github.com/binbankm/My.git
cd My

# Start backend
go run main.go

# Start frontend (new terminal)
cd frontend
npm install
npm run dev
```

### 📦 Building

```bash
# Build all platforms
make build-all

# Build Linux only
make build-linux

# Build frontend only
cd frontend && npm run build

# Build backend only
go build -o serverpanel main.go
```

### 🎯 Roadmap

- [x] Basic system monitoring
- [x] File management
- [x] User authentication
- [ ] Complete Docker integration
- [ ] Database connection management (MySQL/PostgreSQL)
- [ ] Nginx configuration management
- [ ] Scheduled tasks
- [ ] Log viewer
- [ ] Terminal/SSH integration
- [ ] Backup and restore
- [ ] Multi-user permissions
- [ ] WebSocket real-time communication

### 🤝 Contributing

Contributions are welcome! Please see [DEVELOPMENT.md](docs/DEVELOPMENT.md) for how to get started.

### 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

---

<div align="center">

Made with ❤️ by binbankm

⭐ Star this repo if you find it useful!

</div>
