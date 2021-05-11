GOFLAGS			 = -mod=readonly
GOLANGCI_LINT	 = github.com/golangci/golangci-lint/cmd/golangci-lint

up:
	@docker-compose run --rm start_dependencies
	@docker-compose up -d cart-server cart-consumer

down:
	@docker-compose down

build:
	@docker-compose build

validate:
	@printf "=== $(INTEGRATION) === [ validate ]: running golangci-lint & semgrep... "
	@go run  $(GOFLAGS) $(GOLANGCI_LINT) run --verbose
	@if [ -f .semgrep.yml ]; then \
        docker run --rm -v "${PWD}:/src:ro" --workdir /src returntocorp/semgrep -c .semgrep.yml ; \
    else \
    	docker run --rm -v "${PWD}:/src:ro" --workdir /src returntocorp/semgrep -c p/golang ; \
    fi

test:
	@echo ">>> Running tests..."
	@go test -race -v ./... -count=1
	
.PHONY: up down build test
