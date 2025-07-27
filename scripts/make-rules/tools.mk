# ==============================================================================
# 工具相关的 Makefile
#

TOOLS ?= golangci-lint goimports protoc-plugins swagger

.PHONY: tools.verify
tools.verify: $(addprefix tools.verify., $(TOOLS))

.PHONY: tools.install
tools.install: $(addprefix tools.install., $(TOOLS))

.PHONY: tools.install.all
tools.install.all: tools.install.golangci-lint tools.install.goimports tools.install.wire tools.install.protoc-plugins tools.install.swagger

.PHONY: tools.install.%
tools.install.%:
	@echo "===========> Installing $*"
	@$(MAKE) install.$*

.PHONY: tools.verify.%
tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi

.PHONY: install.golangci-lint
install.golangci-lint:
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: install.goimports
install.goimports:
	@$(GO) install golang.org/x/tools/cmd/goimports@latest

.PHONY: install.wire
install.wire:
	@$(GO) install github.com/google/wire/cmd/wire@latest

.PHONY: install.protoc-plugins
install.protoc-plugins:
	@$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: install.swagger
install.swagger:
	@$(GO) install github.com/go-swagger/go-swagger/cmd/swagger@latest
