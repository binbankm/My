# API Documentation

This document describes the REST API endpoints available in ServerPanel.

## Authentication

All API endpoints (except login and register) require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <token>
```

## Base URL

```
http://your-server-ip:8888/api
```

## Authentication Endpoints

### POST /auth/login

Login to the system.

**Request Body:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "admin"
  }
}
```

### POST /auth/register

Register a new user.

**Request Body:**
```json
{
  "username": "newuser",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 2,
    "username": "newuser"
  }
}
```

### GET /auth/me

Get current user information.

**Response:**
```json
{
  "id": 1,
  "username": "admin"
}
```

## System Monitoring Endpoints

### GET /system/info

Get system information.

**Response:**
```json
{
  "hostname": "server-01",
  "platform": "linux",
  "platformVersion": "Ubuntu 22.04",
  "kernelVersion": "5.15.0-91-generic",
  "cpuCores": 4,
  "uptime": 1234567
}
```

### GET /system/stats

Get real-time system statistics.

**Response:**
```json
{
  "cpu": [25.5, 30.2, 15.8, 40.1],
  "memory": {
    "total": 8589934592,
    "used": 4294967296,
    "free": 4294967296,
    "usedPercent": 50.0
  },
  "disk": [
    {
      "path": "/",
      "fstype": "ext4",
      "total": 107374182400,
      "used": 53687091200,
      "free": 53687091200,
      "usedPercent": 50.0
    }
  ],
  "network": {
    "bytesSent": 1073741824,
    "bytesRecv": 2147483648
  }
}
```

## File Management Endpoints

### GET /files

List files and directories.

**Query Parameters:**
- `path` (optional): Directory path to list. Default: `/`

**Response:**
```json
{
  "path": "/home/user",
  "files": [
    {
      "name": "document.txt",
      "size": 1024,
      "mode": "-rw-r--r--",
      "modTime": "2024-01-15T10:30:00Z",
      "isDir": false
    },
    {
      "name": "folder",
      "size": 4096,
      "mode": "drwxr-xr-x",
      "modTime": "2024-01-15T10:30:00Z",
      "isDir": true
    }
  ]
}
```

### GET /files/download

Download a file.

**Query Parameters:**
- `path` (required): File path to download

**Response:**
File content with appropriate Content-Type header

### POST /files/upload

Upload a file.

**Request:**
- Content-Type: `multipart/form-data`
- Form field `file`: File to upload
- Form field `path`: Target directory path

**Response:**
```json
{
  "message": "File uploaded successfully",
  "path": "/uploads/file.txt"
}
```

### POST /files/create

Create a new file or directory.

**Request Body:**
```json
{
  "path": "/home/user/newfile.txt",
  "content": "file content",
  "isDir": false
}
```

**Response:**
```json
{
  "message": "File created successfully"
}
```

### DELETE /files

Delete a file or directory.

**Query Parameters:**
- `path` (required): Path to delete

**Response:**
```json
{
  "message": "File deleted successfully"
}
```

### PUT /files/rename

Rename or move a file.

**Request Body:**
```json
{
  "oldPath": "/home/user/oldname.txt",
  "newPath": "/home/user/newname.txt"
}
```

**Response:**
```json
{
  "message": "File renamed successfully"
}
```

## Docker Container Management Endpoints

### GET /containers

List all Docker containers.

**Response:**
```json
{
  "containers": [
    {
      "id": "abc123",
      "name": "nginx",
      "image": "nginx:latest",
      "status": "running",
      "ports": ["80:80", "443:443"],
      "created": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### POST /containers/start

Start a container.

**Request Body:**
```json
{
  "id": "abc123"
}
```

**Response:**
```json
{
  "message": "Container started successfully"
}
```

### POST /containers/stop

Stop a container.

**Request Body:**
```json
{
  "id": "abc123"
}
```

**Response:**
```json
{
  "message": "Container stopped successfully"
}
```

### DELETE /containers

Remove a container.

**Query Parameters:**
- `id` (required): Container ID

**Response:**
```json
{
  "message": "Container removed successfully"
}
```

## Database Management Endpoints

### GET /databases

List all database connections.

**Response:**
```json
{
  "databases": [
    {
      "id": 1,
      "name": "MySQL Production",
      "type": "mysql",
      "host": "localhost",
      "port": 3306,
      "status": "connected"
    }
  ]
}
```

### POST /databases

Create a new database connection.

**Request Body:**
```json
{
  "name": "MySQL Production",
  "type": "mysql",
  "host": "localhost",
  "port": 3306,
  "username": "root",
  "password": "password",
  "database": "mydb"
}
```

**Response:**
```json
{
  "id": 1,
  "message": "Database connection created successfully"
}
```

### POST /databases/test

Test a database connection.

**Request Body:**
```json
{
  "type": "mysql",
  "host": "localhost",
  "port": 3306,
  "username": "root",
  "password": "password",
  "database": "mydb"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Connection successful"
}
```

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "Invalid request parameters"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

## Rate Limiting

The API does not currently implement rate limiting, but this may be added in future versions.

## Versioning

The current API version is v1. The version is not included in the URL path. Future versions will be clearly documented with migration guides.

## Docker Management Endpoints

### GET /docker/containers
List all Docker containers.

**Response:** Array of container objects

### GET /docker/containers/:id
Get details of a specific container.

### GET /docker/containers/:id/logs
Get logs from a container.

### GET /docker/containers/:id/stats
Get real-time statistics from a container.

### POST /docker/containers/:id/start
Start a container.

### POST /docker/containers/:id/stop
Stop a container.

### POST /docker/containers/:id/restart
Restart a container.

### DELETE /docker/containers/:id
Delete a container.

### GET /docker/images
List all Docker images.

### DELETE /docker/images/:id
Delete a Docker image.

## Database Management Endpoints

### GET /database
List all database connections.

### POST /database
Create a new database connection.

**Request Body:**
```json
{
  "name": "mydb",
  "type": "mysql",
  "host": "localhost",
  "port": 3306,
  "username": "root",
  "password": "password",
  "database": "mydb"
}
```

### GET /database/:id
Get a specific database connection.

### DELETE /database/:id
Delete a database connection.

### POST /database/:id/test
Test a database connection.

### POST /database/:id/query
Execute a SQL query.

**Request Body:**
```json
{
  "query": "SELECT * FROM users"
}
```

### GET /database/:id/tables
List all tables in a database.

## Cron Job Management Endpoints

### GET /cron
List all cron jobs.

### POST /cron
Create a new cron job.

**Request Body:**
```json
{
  "schedule": "0 2 * * *",
  "command": "/path/to/script.sh",
  "comment": "Daily backup"
}
```

### GET /cron/:id
Get a specific cron job.

### PUT /cron/:id
Update a cron job.

### DELETE /cron/:id
Delete a cron job.

## Log Viewer Endpoints

### GET /logs/files
List all available log files.

### GET /logs/read
Read content from a log file.

**Query Parameters:**
- `path`: Path to log file (required)
- `lines`: Number of lines to read (default: 100)
- `tail`: Read from end of file (default: true)
- `filter`: Filter by keyword

### GET /logs/search
Search across multiple log files.

**Query Parameters:**
- `query`: Search query (required)
- `dir`: Directory to search (default: /var/log)

### GET /logs/system
Get system logs using journalctl.

**Query Parameters:**
- `lines`: Number of lines (default: 100)
- `unit`: Systemd unit name
- `since`: Time range

### GET /logs/download
Download a log file.

**Query Parameters:**
- `path`: Path to log file (required)

### POST /logs/clear
Clear a log file.

**Request Body:**
```json
{
  "path": "/var/log/myapp.log"
}
```

### GET /logs/stats
Get log file statistics.

## Nginx Management Endpoints

### GET /nginx/sites
List all Nginx sites.

### GET /nginx/sites/:name
Get a specific Nginx site configuration.

### POST /nginx/sites
Create a new Nginx site.

**Request Body:**
```json
{
  "name": "example.com",
  "serverName": "example.com",
  "port": "80",
  "root": "/var/www/example",
  "content": "server { ... }"
}
```

### PUT /nginx/sites/:name
Update an Nginx site configuration.

### DELETE /nginx/sites/:name
Delete an Nginx site.

### POST /nginx/sites/:name/enable
Enable an Nginx site.

### POST /nginx/sites/:name/disable
Disable an Nginx site.

### POST /nginx/test
Test Nginx configuration.

### POST /nginx/reload
Reload Nginx.

### GET /nginx/status
Get Nginx status.

## Backup and Restore Endpoints

### GET /backup
List all backups.

### POST /backup
Create a new backup.

**Request Body:**
```json
{
  "type": "file",
  "name": "mybackup",
  "source": "/path/to/data",
  "description": "Important data backup"
}
```

### GET /backup/:id/download
Download a backup file.

### DELETE /backup/:id
Delete a backup.

### POST /backup/:id/restore
Restore from a backup.

**Request Body:**
```json
{
  "destination": "/path/to/restore"
}
```

### GET /backup/stats
Get backup statistics.

## User Management Endpoints

### GET /users
List all users.

### GET /users/:id
Get a specific user.

### POST /users
Create a new user.

**Request Body:**
```json
{
  "username": "newuser",
  "password": "password",
  "email": "user@example.com",
  "isAdmin": false,
  "roleId": 2
}
```

### PUT /users/:id
Update a user.

### DELETE /users/:id
Delete a user.

## Role Management Endpoints

### GET /roles
List all roles.

### GET /roles/:id
Get a specific role.

### POST /roles
Create a new role.

**Request Body:**
```json
{
  "name": "developer",
  "description": "Developer role",
  "permissionIds": [1, 3, 5, 8]
}
```

### PUT /roles/:id
Update a role.

### DELETE /roles/:id
Delete a role.

## Permission Endpoints

### GET /permissions
List all permissions.

## WebSocket Endpoint

### GET /ws
Connect to WebSocket for real-time updates.

WebSocket messages are in JSON format:
```json
{
  "type": "system_stats",
  "data": {
    "cpu": {...},
    "memory": {...}
  },
  "timestamp": 1234567890
}
```

