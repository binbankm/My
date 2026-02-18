# 部署指南 / Deployment Guide

## 系统要求 / System Requirements

- Linux (Ubuntu 20.04+, CentOS 7+, Debian 10+ recommended)
- 1GB+ RAM
- 10GB+ disk space
- Root access or sudo privileges

## 快速安装 / Quick Installation

### 方法 1: 一键远程安装（推荐）/ Method 1: One-line Remote Installation (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/binbankm/My/main/scripts/remote-install.sh | sudo bash
```

或者使用 wget / Or using wget:

```bash
wget -qO- https://raw.githubusercontent.com/binbankm/My/main/scripts/remote-install.sh | sudo bash
```

### 方法 2: 使用安装脚本 / Method 2: Using Install Script

```bash
# 下载最新版本 / Download latest release
wget https://github.com/binbankm/My/releases/latest/download/serverpanel-linux-amd64.tar.gz

# 解压 / Extract
tar -zxvf serverpanel-linux-amd64.tar.gz

# 运行安装脚本 / Run install script
sudo ./install.sh
```

### 方法 3: 手动安装 / Method 3: Manual Installation

```bash
# 1. 创建目录 / Create directory
sudo mkdir -p /opt/serverpanel

# 2. 复制二进制文件 / Copy binary
sudo cp serverpanel-linux-amd64 /opt/serverpanel/serverpanel
sudo chmod +x /opt/serverpanel/serverpanel

# 3. 创建systemd服务 / Create systemd service
sudo nano /etc/systemd/system/serverpanel.service
```

添加以下内容 / Add the following content:

```ini
[Unit]
Description=ServerPanel - Linux Server Management
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/serverpanel
ExecStart=/opt/serverpanel/serverpanel
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

```bash
# 4. 启动服务 / Start service
sudo systemctl daemon-reload
sudo systemctl enable serverpanel
sudo systemctl start serverpanel
```

## 访问面板 / Access Panel

安装完成后，在浏览器中访问:
After installation, access in browser:

```
http://YOUR_SERVER_IP:8888
```

默认登录凭证 / Default credentials:
- 用户名 / Username: `admin`
- 密码 / Password: `admin123`

**重要**: 首次登录后请立即修改密码！
**Important**: Change the password immediately after first login!

## 配置 / Configuration

### 更改端口 / Change Port

编辑服务文件 / Edit service file:

```bash
sudo nano /etc/systemd/system/serverpanel.service
```

在 ExecStart 行添加环境变量:
Add environment variable to ExecStart line:

```ini
Environment="PORT=9000"
ExecStart=/opt/serverpanel/serverpanel
```

重启服务 / Restart service:

```bash
sudo systemctl daemon-reload
sudo systemctl restart serverpanel
```

### 防火墙设置 / Firewall Setup

#### Ubuntu/Debian (UFW)

```bash
sudo ufw allow 8888/tcp
sudo ufw reload
```

#### CentOS/RHEL (firewalld)

```bash
sudo firewall-cmd --permanent --add-port=8888/tcp
sudo firewall-cmd --reload
```

## 常用命令 / Common Commands

```bash
# 查看状态 / Check status
sudo systemctl status serverpanel

# 启动 / Start
sudo systemctl start serverpanel

# 停止 / Stop
sudo systemctl stop serverpanel

# 重启 / Restart
sudo systemctl restart serverpanel

# 查看日志 / View logs
sudo journalctl -u serverpanel -f

# 查看最近100行日志 / View last 100 log lines
sudo journalctl -u serverpanel -n 100
```

## 卸载 / Uninstall

**使用远程卸载脚本 / Using remote uninstall script:**

```bash
wget -qO- https://raw.githubusercontent.com/binbankm/My/main/scripts/uninstall.sh | sudo bash
```

**或使用本地卸载脚本 / Or using local uninstall script:**

```bash
# 使用卸载脚本 / Using uninstall script
sudo ./uninstall.sh

# 或手动卸载 / Or manually uninstall
sudo systemctl stop serverpanel
sudo systemctl disable serverpanel
sudo rm /etc/systemd/system/serverpanel.service
sudo rm -rf /opt/serverpanel
sudo systemctl daemon-reload
```

## 升级 / Upgrade

```bash
# 1. 停止服务 / Stop service
sudo systemctl stop serverpanel

# 2. 备份数据库 / Backup database
sudo cp /opt/serverpanel/serverpanel.db /opt/serverpanel/serverpanel.db.backup

# 3. 下载新版本 / Download new version
wget https://github.com/binbankm/My/releases/latest/download/serverpanel-linux-amd64.tar.gz

# 4. 解压并替换 / Extract and replace
tar -zxvf serverpanel-linux-amd64.tar.gz
sudo cp serverpanel-linux-amd64 /opt/serverpanel/serverpanel

# 5. 启动服务 / Start service
sudo systemctl start serverpanel
```

## 故障排查 / Troubleshooting

### 服务无法启动 / Service won't start

```bash
# 检查端口是否被占用 / Check if port is in use
sudo netstat -tlnp | grep 8888

# 检查日志 / Check logs
sudo journalctl -u serverpanel -n 50 --no-pager
```

### 无法访问面板 / Can't access panel

1. 检查服务状态 / Check service status: `sudo systemctl status serverpanel`
2. 检查防火墙 / Check firewall settings
3. 确认服务器IP地址 / Confirm server IP address
4. 检查是否使用了代理或云服务商的安全组规则 / Check if using proxy or cloud security groups

### 数据库问题 / Database issues

```bash
# 查看数据库文件 / View database file
ls -lh /opt/serverpanel/serverpanel.db

# 从备份恢复 / Restore from backup
sudo cp /opt/serverpanel/serverpanel.db.backup /opt/serverpanel/serverpanel.db
sudo systemctl restart serverpanel
```

## 安全建议 / Security Recommendations

1. **修改默认密码** / Change default password immediately
2. **使用HTTPS** / Use HTTPS (配置Nginx反向代理 / Configure Nginx reverse proxy)
3. **限制访问IP** / Restrict access by IP (使用防火墙规则 / Use firewall rules)
4. **定期备份** / Regular backups (数据库文件 / Database file)
5. **保持更新** / Keep updated (定期检查新版本 / Check for updates regularly)

## 性能优化 / Performance Optimization

### 对于大流量场景 / For high traffic scenarios

考虑使用Nginx作为反向代理:
Consider using Nginx as reverse proxy:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8888;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 支持 / Support

- GitHub Issues: https://github.com/binbankm/My/issues
- 文档 / Documentation: https://github.com/binbankm/My/wiki

## 许可证 / License

MIT License - 详见 LICENSE 文件 / See LICENSE file
