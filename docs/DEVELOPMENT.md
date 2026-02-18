# Development Guide

## Prerequisites

- Go 1.21 or later
- Node.js 18 or later
- npm or pnpm

## Project Structure

```
.
├── main.go                 # Main application entry point
├── internal/              # Private application code
│   ├── api/              # API handlers
│   ├── middleware/       # Middleware functions
│   ├── model/            # Database models
│   ├── service/          # Business logic
│   └── util/             # Utility functions
├── frontend/             # Vue.js frontend
│   ├── src/
│   │   ├── components/  # Vue components
│   │   ├── views/       # Page views
│   │   ├── router/      # Vue Router config
│   │   ├── store/       # Pinia stores
│   │   └── api/         # API client
│   └── dist/            # Built frontend (generated)
├── scripts/             # Installation scripts
├── docs/                # Documentation
└── .github/workflows/   # CI/CD workflows
```

## Development Setup

### Backend Development

```bash
# Install Go dependencies
go mod download

# Run backend (without frontend)
go run main.go
```

The backend API will be available at `http://localhost:8888`

### Frontend Development

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The frontend dev server will run at `http://localhost:3000` with hot reload.

### Full Stack Development

Terminal 1 (Backend):
```bash
go run main.go
```

Terminal 2 (Frontend):
```bash
cd frontend
npm run dev
```

Access the application at `http://localhost:3000` (frontend dev server with API proxy).

## Building

### Build Frontend Only

```bash
cd frontend
npm run build
```

### Build Backend Only

```bash
go build -o serverpanel main.go
```

### Build Both (Production)

```bash
# Using Makefile
make build-backend  # Builds both frontend and backend
```

## Testing

### Backend Tests

```bash
go test -v ./...
```

### Frontend Tests

```bash
cd frontend
npm run test  # If tests are configured
```

## Code Style

### Go

Follow standard Go conventions:
- Use `gofmt` for formatting
- Follow effective Go guidelines
- Keep functions small and focused

### JavaScript/Vue

- Use ESLint configuration
- Follow Vue.js style guide
- Use Composition API for new components

## Adding New Features

### Backend API Endpoint

1. Add handler in `internal/api/`
2. Add route in `main.go`
3. Add middleware if needed
4. Update API documentation

### Frontend Page

1. Create view component in `frontend/src/views/`
2. Add route in `frontend/src/router/index.js`
3. Add API calls in `frontend/src/api/index.js`
4. Update navigation if needed

## Database

The application uses SQLite by default. The database file is `serverpanel.db`.

### Migrations

Migrations are handled automatically by GORM's AutoMigrate feature.

### Models

Add new models in `internal/model/model.go` and they will be auto-migrated on startup.

## Debugging

### Backend

Use your favorite Go debugger or add print statements:

```go
log.Printf("Debug: %+v", variable)
```

### Frontend

Use browser DevTools:
- Console for logs
- Network tab for API calls
- Vue DevTools extension for Vue-specific debugging

## Environment Variables

- `PORT`: Server port (default: 8888)
- `GIN_MODE`: Gin mode (`debug` or `release`)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Release Process

Releases are automated via GitHub Actions:

1. Merge changes to `main` branch
2. Auto-tag workflow creates a new version tag
3. Release workflow builds binaries for all platforms
4. Artifacts are attached to the GitHub release

Manual tagging:
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```
