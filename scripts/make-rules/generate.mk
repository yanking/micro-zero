# ==============================================================================
# 用来进行代码生成的 Makefile
#

.PHONY: gen.protoc
gen.protoc: # 生成gRPC相关文件
	@echo "===========> Generate protobuf files"
	@#buf dep update $(APIROOT)
	@buf format -w $(APIROOT)
	@buf lint $(APIROOT)
	@buf generate $(APIROOT)

.PHONY: gen.generate
gen.generate:
	@GOWORK=off go generate ./...

