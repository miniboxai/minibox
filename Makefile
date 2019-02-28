
# include .env_make

PROJECTNAME=$(shell basename "$(PWD)")
PROJECTCLI=mini
PROJECTGEN=minisca

ADDR=14000
# Go related variables.
GOBASE=$(shell pwd)
# GOPATH=$(GOBASE)/vendor:$(GOBASE)
GOBIN=$(GOBASE)/bin
# GOFILES=$(wildcard *.go)
GOFILES=cmd/apiserver/*.go
PARENTDIR=$(shell dirname "$(GOBASE)")

ifeq ($(PROD),)
GOFLAGS=-gcflags="-l=4 "
GOOPTS=-tags="server" $(GOFLAGS)
else
GOFLAGS=-gcflags "all=-trimpath=$(PARENTDIR)" -asmflags "all=-trimpath=$(PARENTDIR)"
GOOPTS=-tags="release server" $(GOFLAGS)
endif

# GOFLAGS=-gcflags "all=-trimpath=$(PARENTDIR)" -asmflags "all=-trimpath=$(PARENTDIR)"

GOCLI=cmd/mini/*.go
# PBFILES=$(wildcard pkg/api/v1/**/*.proto)

GE_INCLUDES = -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
 -I=$GOPATH/src/github.com/gogo/protobuf/protobuf

PROTO_FILES = $(shell find pkg/api/v1 -type f -name '*.proto')
PBGO_FILES = $(patsubst pkg/api/v1/%.proto, pkg/api/v1/%.pb.go, $(PROTO_FILES))
GATEWAY_FILES = $(patsubst pkg/api/v1/%.proto, pkg/api/v1/%.gw.go, $(PROTO_FILES))
SWAGGER_FILES = $(patsubst pkg/api/v1/%.proto, pkg/api/v1/%.swagger.json, $(PROTO_FILES))

# Redirect error output to a file, so we can show it in development mode.
STDERR=/tmp/.$(PROJECTNAME)-stderr.txt

# PID file will keep the process id of the server
PID=/tmp/.$(PROJECTNAME).pid

LOG=/tmp/$(PROJECTNAME).log

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## install: Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
install: go-install


## start: Start in development mode. Auto-starts when code changes.
start:
	bash -c "trap 'make stop' EXIT; $(MAKE) compile start-server watch run='make compile start-server'"

watch-cli:
	bash -c "trap 'make stop' EXIT; $(MAKE) cli watch run='make cli'"

wathc-predata: 
	bash -c "trap 'make stop' EXIT; $(MAKE) cli watch run='make predata'"

## stop: Stop development mode.
stop: stop-server

start-server: stop-server
	@echo "  >  $(PROJECTNAME) is available at $(ADDR)"
	@-$(GOBIN)/$(PROJECTNAME) >> $(LOG) 2>&1 & echo $$! > $(PID)
	@cat $(PID) | sed "/^/s/^/  \>  PID: /"

stop-server:
	@-touch $(PID)
	@-kill `cat $(PID)` 2> /dev/null || true
	@-rm $(PID)

## watch: Run given command when code changes. e.g; make watch run="echo 'hey'"
watch:
	@yolo -i . -e db -e vendor -e bin -c "$(run)" -a localhost:8081

restart-server: stop-server start-server

## compile: Compile the binary.
compile:
	@-touch $(STDERR)
	@-rm $(STDERR)
	@-$(MAKE) -s go-compile 2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/.*/\nError:\n/'  | sed 's/make\[.*/ /' | sed "/^/s/^/     /" 1>&2

## exec: Run given command, wrapped with custom GOPATH. e.g; make exec run="go test ./..."
exec:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) $(run)

## clean: Clean build files. Runs `go clean` internally.
clean:
	@(MAKEFILE) go-clean

go-compile: go-clean go-get go-build

go-compile-cli: go-clean go-get go-build-cli

go-compile-predata: go-clean go-get go-build-predata

go-build:
	@echo "  >  Building binary..."
ifdef PROD
	@echo "  >  \033[0;32mRelease Mode\033[0m"
endif
	@go build $(GOOPTS) -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

go-build-cli:
	@echo "  >  Building binary cli..."
	@go build $(GOFLAGS) -o $(GOBIN)/$(PROJECTCLI) $(GOCLI)
	# @echo "  >  Building binary minisca..."
	# @go build $(GOFLAGS) -o $(GOBIN)/$(PROJECTGEN) cmd/minisca/*.go

go-build-predata:
	@echo "  >  Building binary predata..."
	@go build $(GOFLAGS) -o $(GOBIN)/predata cmd/predata/*.go

go-generate:
	@echo "  >  Generating dependency files..."
	@go generate $(GOOPTS) $(generate)

go-get:
	@echo "  >  Checking if there is any missing dependencies..."

go-install:
	@go install $(GOOPTS) $(GOFILES)
	@GO111MODULE=off go get github.com/azer/yolo
	@GO111MODULE=off go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	@GO111MODULE=off go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	@GO111MODULE=off go get github.com/golang/protobuf/protoc-gen-go
	@GO111MODULE=off go get github.com/elazarl/go-bindata-assetfs/...
	@GO111MODULE=off go get github.com/jteeuwen/go-bindata/...
	@GO111MODULE=off go get github.com/gogo/protobuf/protoc-gen-gofast
	@GO111MODULE=off go get github.com/gogo/protobuf/proto
	@GO111MODULE=off go get github.com/gogo/protobuf/jsonpb
	@GO111MODULE=off go get github.com/gogo/protobuf/protoc-gen-gogo
	@GO111MODULE=off go get github.com/gogo/protobuf/gogoproto/...
	@GO111MODULE=off go get github.com/gogo/googleapis/...
	@GO111MODULE=off go get github.com/mwitkow/go-proto-validators
	@GO111MODULE=off go get github.com/golang/mock/gomock
	@GO111MODULE=off go install github.com/golang/mock/mockgen

go-clean:
	@echo "  >  Cleaning build cache"
	# @go clean

go-protoc: $(PBGO_FILES)

# go-protoc: $(PBGO_FILES) $(GATEWAY_FILES) $(SWAGGER_FILES)

# pkg/api/v1/%.pb.go: pkg/api/v1/%.proto
# 	$(eval dir=$(shell dirname $?))
# 	@echo convert "$<" to "$@"
# 	# @protoc -I $(dir) $(GE_INCLUDES) --go_out=plugins=grpc:$(dir) $? 
# 	@protoc -I $(dir) $(GE_INCLUDES) --gofast_out=plugins=grpc:$(dir) $? 

# pkg/api/v1/%.gw.go: pkg/api/v1/%.proto
# 	$(eval dir=$(shell dirname $?))
# 	@echo convert gateway "$<" to "$@"
# 	protoc -I $(dir) $(GE_INCLUDES) --grpc-gateway_out=logtostderr=true:$(dir)  $?

# pkg/api/v1/%.swagger.json: pkg/api/v1/%.proto
# 	$(eval dir=$(shell dirname $?))
# 	@echo convert swagger json "$<" to "$@"
# 	@protoc -I $(dir) $(GE_INCLUDES) --swagger_out=logtostderr=true:$(dir)  $?

clean-protoc: $(PBGO_FILES)
	@echo rm $<

pkg/api/v1/%.pb.go: pkg/api/v1/%.proto
	$(eval dir=$(shell dirname $?))
	@echo convert gateway "$<" to "$@"
	@GO111MODULE=off protoc \
			-I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/ \
			-I $(GOPATH)/src/github.com/gogo/googleapis/ \
			-I $(GOPATH)/src \
			-I pkg/api/v1/types \
			-I $(dir) \
			-I vendor/ \
			--gogo_out=plugins=grpc,\
	Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
	Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,\
	Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,\
	Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,\
	Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,\
	Mtypes.proto=minibox.ai/pkg/api/v1/types:\
	$(dir) \
			--grpc-gateway_out=\
	Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
	Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,\
	Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,\
	Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,\
	Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,\
	Mtypes.proto=minibox.ai/pkg/api/v1/types:\
	$(dir) \
			--swagger_out=third_party/OpenAPI/ \
			--govalidators_out=gogoimport=true,\
	Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
	Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,\
	Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,\
	Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,\
	Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,\
	Mtypes.proto=minibox.ai/pkg/api/v1/types:\
	$(dir) \
			$?

hack: hack-protoc

hack-protoc: 
	find . -name "*.gw.go" -exec sed -i '' 's/ListObjectsRequest_DatasetId/types.ListObjectsRequest_DatasetId/g' {} +

cli: go-compile-cli

predata: go-compile-predata

generate_predata:
	@make go-generate generate="./pkg/predata"

generate:
	@make go-generate generate="./pkg/api/..."
	@make go-generate generate="./pkg/server/..."

swagger:
	@docker run -d -p 80:8080 -e API_URL=http://generator.swagger.io/api/swagger.json swaggerapi/swagger-ui
	@open http://localhost

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
