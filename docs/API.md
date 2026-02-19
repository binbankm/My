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
