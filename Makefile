.PHONY: all build-frontend build-backend build-linux build-darwin build-windows build-all clean

all: build-all

build-frontend:
	@echo "Building frontend..."
	cd frontend && npm install && npm run build

build-backend: build-frontend
	@echo "Building backend..."
	go mod download
	go build -o bin/serverpanel main.go

build-linux: build-frontend
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o bin/serverpanel-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o bin/serverpanel-linux-arm64 main.go

build-darwin: build-frontend
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build -o bin/serverpanel-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/serverpanel-darwin-arm64 main.go

build-windows: build-frontend
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o bin/serverpanel-windows-amd64.exe main.go

build-all: build-linux build-darwin build-windows
	@echo "All builds complete!"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf frontend/dist/
	rm -rf frontend/node_modules/

test:
	go test -v ./...

run: build-backend
	./bin/serverpanel
