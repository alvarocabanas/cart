up:
	@docker-compose run --rm start_dependencies
	@docker-compose up -d cart-server cart-consumer

down:
	@docker-compose down

build:
	@docker-compose build

test:
	@echo ">>> Running tests..."
	@go test -race -v ./... -count=1