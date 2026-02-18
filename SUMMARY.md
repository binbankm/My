# ServerPanel - Project Summary

## Overview

ServerPanel is a production-ready Linux server management panel built with Go and Vue.js, similar to 1Panel. It provides a modern web interface for managing Linux servers with real-time monitoring, file management, and extensible architecture.

## Key Features Implemented

### ✅ Backend (Go)
- **RESTful API** - Built with Gin framework
- **Authentication** - JWT tokens with bcrypt password hashing
- **System Monitoring** - Real-time CPU, memory, disk, and network monitoring using gopsutil
- **File Management** - Complete file operations with path validation and security
- **User Management** - SQLite database with GORM
- **Security** - Environment-based configuration, CORS protection, path sandboxing
- **Single Binary** - Embedded frontend assets using go:embed

### ✅ Frontend (Vue 3)
- **Modern UI** - Built with Tailwind CSS and shadcn/ui components
- **Real-time Dashboard** - Live system statistics with auto-refresh
- **File Manager** - Directory navigation with upload/download support
- **Container Management** - Docker interface (API ready for extension)
- **Authentication** - JWT-based login with secure token storage
- **Responsive Design** - Works on desktop and mobile devices

### ✅ DevOps & Deployment
- **Multi-platform Builds** - Linux (amd64/arm64), macOS (amd64/arm64), Windows (amd64)
- **GitHub Actions** - Automated tagging, building, and releasing
- **Systemd Integration** - Installation and service management scripts
- **Single Binary Distribution** - ~18MB executable with embedded frontend
- **Documentation** - Comprehensive guides in English and Chinese

## Architecture

```
┌─────────────────────────────────────────┐
│          Browser (Frontend)            │
│     Vue 3 + Tailwind + shadcn/ui       │
└──────────────┬──────────────────────────┘
               │ HTTP/REST API
               │
┌──────────────▼──────────────────────────┐
│         Go Backend (Gin)               │
│  ┌────────────────────────────────┐    │
│  │  API Layer                     │    │
│  │  - Auth, System, Files, etc.   │    │
│  └────────────────────────────────┘    │
│  ┌────────────────────────────────┐    │
│  │  Middleware                    │    │
│  │  - JWT Auth, CORS, Logging     │    │
│  └────────────────────────────────┘    │
│  ┌────────────────────────────────┐    │
│  │  Services                      │    │
│  │  - System Info, File Ops       │    │
│  └────────────────────────────────┘    │
│  ┌────────────────────────────────┐    │
│  │  Database (SQLite + GORM)      │    │
│  └────────────────────────────────┘    │
└─────────────────────────────────────────┘
```

## Security Features

1. **JWT Authentication** - Secure token-based auth
2. **Password Hashing** - bcrypt with cost factor 10
3. **Environment Configuration** - Secrets loaded from environment variables
4. **CORS Protection** - Configurable origin restrictions
5. **Path Validation** - File operations sandboxed to allowed directories
6. **Input Validation** - Request data validation and sanitization

## File Structure

```
serverpanel/
├── main.go                    # Application entry point
├── internal/
│   ├── api/                   # HTTP handlers
│   │   ├── auth.go           # Authentication endpoints
│   │   ├── system.go         # System monitoring
│   │   ├── files.go          # File management
│   │   ├── docker.go         # Container management
│   │   ├── database.go       # Database operations
│   │   └── settings.go       # Settings management
│   ├── middleware/
│   │   └── auth.go           # CORS and JWT middleware
│   ├── model/
│   │   └── model.go          # Database models
│   └── util/
│       └── auth.go           # Auth utilities
├── frontend/
│   ├── src/
│   │   ├── views/            # Page components
│   │   ├── components/       # Reusable components
│   │   ├── router/           # Vue Router
│   │   ├── store/            # Pinia stores
│   │   └── api/              # API client
│   └── dist/                 # Built assets (embedded)
├── scripts/
│   ├── install.sh            # Installation script
│   └── uninstall.sh          # Uninstallation script
├── docs/
│   ├── DEPLOYMENT.md         # Deployment guide
│   └── DEVELOPMENT.md        # Development guide
└── .github/workflows/        # CI/CD workflows
```

## Configuration

Environment variables (see `.env.example`):

- `PORT` - Server port (default: 8888)
- `GIN_MODE` - Gin mode (release/debug)
- `JWT_SECRET` - JWT signing secret (required in production)
- `CORS_ORIGINS` - Allowed CORS origins
- `FILE_MANAGER_BASE_PATH` - Base path for file operations

## Build Information

- **Go Version**: 1.21+
- **Node Version**: 18+
- **Binary Size**: ~18MB (includes frontend)
- **Database**: SQLite (embedded)
- **Supported Platforms**: Linux, macOS, Windows

## API Endpoints

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout
- `GET /api/auth/info` - Get user info

### System Monitoring
- `GET /api/system/info` - System information
- `GET /api/system/stats` - Real-time statistics

### File Management
- `GET /api/files` - List files
- `POST /api/files` - Create file/directory
- `PUT /api/files` - Update file content
- `DELETE /api/files` - Delete file
- `GET /api/files/download` - Download file
- `POST /api/files/upload` - Upload file

### Docker (API Ready)
- `GET /api/docker/containers` - List containers
- `POST /api/docker/containers/:id/start` - Start container
- `POST /api/docker/containers/:id/stop` - Stop container
- etc.

## Future Enhancements

Potential features for future versions:

- [ ] Complete Docker integration with real container management
- [ ] MySQL/PostgreSQL connection management
- [ ] Nginx configuration editor
- [ ] Cron job scheduler
- [ ] Log viewer with search and filtering
- [ ] Terminal/SSH integration (xterm.js)
- [ ] Backup and restore functionality
- [ ] Multi-user with role-based permissions
- [ ] WebSocket for real-time updates
- [ ] Mobile app (React Native)
- [ ] Prometheus metrics export
- [ ] Email/notification system

## Deployment Recommendations

### Production Checklist

- [x] Change default admin password
- [x] Set JWT_SECRET environment variable
- [x] Configure CORS_ORIGINS
- [ ] Setup HTTPS with reverse proxy (Nginx)
- [ ] Configure firewall rules
- [ ] Setup regular backups
- [ ] Enable system logs rotation
- [ ] Monitor resource usage
- [ ] Keep system updated

### Reverse Proxy Example (Nginx)

```nginx
server {
    listen 80;
    server_name panel.example.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name panel.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8888;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## License

MIT License - See LICENSE file

## Credits

- Inspired by 1Panel
- Built with Go, Vue.js, and modern web technologies
- Community contributions welcome

## Support

- GitHub Issues: https://github.com/binbankm/My/issues
- Documentation: https://github.com/binbankm/My/tree/main/docs

---

**Version**: 1.0.0  
**Status**: Production Ready  
**Last Updated**: 2024
