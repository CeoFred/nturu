# Makefile for generating interfaces from a specified protoc definition

# Set the path to the protoc compiler and plugins
PROTOC := protoc
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(GOPATH)/bin/protoc-gen-go-grpc

# Set the source directory and proto file
PROTO_SRC_DIR := ./shared/grpc
PROTO_FILE := user.proto

# Set the output directory for generated files
OUT_DIR := .

# Generate Go code from the proto file
generate: clean
	cd $(PROTO_SRC_DIR)
	@echo "Generating Go code..."
	$(PROTOC) --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	    --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	    $(PROTO_SRC_DIR)/$(PROTO_FILE)
	@echo "Go code generation complete."
	@echo "Running go mod tidy..."
	go mod tidy
	@echo "go mod tidy complete."

# Create the output directory
$(OUT_DIR):
	mkdir -p $(OUT_DIR)

# Clean the generated files
clean:
	@echo "Cleaning generated files..."
	rm -rf $(PROTO_SRC_DIR)/user.pb.go
	rm -rf $(PROTO_SRC_DIR)/user_grpc.pb.go
	@echo "Clean complete."


.PHONY: generate clean
