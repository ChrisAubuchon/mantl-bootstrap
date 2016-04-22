TEST?=./...
NAME = $(shell awk -F\" '/^const Name/ { print $$2 }' main.go)
VERSION = $(shell awk -F\" '/^const Version/ { print $$2 }' main.go)
DEPS = $(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

all: deps build

deps:
	go get -d -v ./...
	echo $(DEPS) | xargs -n1 go get -d

updatedeps:
	go get -u -v ./...
	echo $(DEPS) | xargs -n1 go get -d

build: deps
	@mkdir -p bin/
	go build -o bin/$(NAME)

test: deps
	go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4
	go vet $(TEST)

xcompile: deps test
	@rm -rf build/
	@mkdir -p build
	gox \
		-os="darwin" \
		-os="linux" \
		-output="build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)"

package: xcompile
	$(eval FILES := $(shell ls build))
	@mkdir -p build/zip
	for f in $(FILES); do \
		(cd $(shell pwd)/build/$$f && zip -q ../zip/$$f.zip mantl-bootstrap); \
		echo $$f; \
	done

pb: 
	protoc -I ./proto ./proto/*.proto --go_out=plugins=grpc:proto

.PHONY: all deps updatedeps build test xcompile package
