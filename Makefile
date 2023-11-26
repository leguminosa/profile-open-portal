

.PHONY: clean all jwt init generate generate_mocks

all: build/main

build/main: cmd/main.go generated
	@echo "Building..."
	go build -o $@ $<

clean:
	rm -rf generated

init: generate jwt
	go mod tidy
	go mod vendor

jwt:
	@echo "Generating rsa..."
	openssl genrsa -out jwt.pem 4096
	openssl rsa -in jwt.pem -pubout -outform PEM -out jwt_pub.pem

test:
	go test -timeout 30s -short -count=1 -race -cover -coverprofile coverage.out -v ./...
	@go tool cover -func coverage.out

generate: generated generate_mocks

generated: api.yml
	@echo "Generating files..."
	mkdir generated || true
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go

INTERFACES_GO_FILES := $(shell find module repository tools -name "interfaces.go")
INTERFACES_GEN_GO_FILES := $(INTERFACES_GO_FILES:%.go=%.mock.gen.go)

generate_mocks: $(INTERFACES_GEN_GO_FILES)
$(INTERFACES_GEN_GO_FILES): %.mock.gen.go: %.go
	@echo "Generating mocks $@ for $<"
	mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))