SSH_KEY_PATH=/home/saimon/.ssh/id_ed25519_no_passphrase
MIN_COVERAGE ?= 0.0

.PHONY: build test coverage

test:
	@echo "🔍 Running tests..."
	go test -v -coverprofile=coverage.out ./...

coverage: test
	@echo "📊 Checking coverage..."
	@coverage=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//') && \
	echo "Total coverage: $${coverage}%" && \
	result=$$(echo "$$coverage >= $(MIN_COVERAGE)" | bc -l) && \
	if [ "$$result" -ne 1 ]; then \
		echo "❌ Coverage below $(MIN_COVERAGE)%!"; \
		exit 1; \
	else \
		echo "✅ Coverage is acceptable."; \
	fi
	
build: coverage
	@echo "🚀 Building Docker image..."
	@SSH_KEY=$(shell cat $(SSH_KEY_PATH) | base64 -w 0) && \
	docker build --build-arg SSH_KEY="$${SSH_KEY}" -t ghcr.io/inter-hubly/pulse:development . && \
	docker push ghcr.io/inter-hubly/pulse:development
