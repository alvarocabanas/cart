up:
	@docker-compose up -d cart

down:
	@docker-compose down

build:
	@docker-compose build

test:
	@echo ">>> Running tests..."
	@go test -race -v ./... -count=1