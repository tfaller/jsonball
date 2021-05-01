today:=$(shell date +%Y-%m-%d)

.DEFAULT_GOAL := all

build/cli/%:
	go build -o $@ $?

cli: build/cli/changes \
	build/cli/gtedoc \
	build/cli/handlernewdoc \
	build/cli/listendocs \
	build/cli/postdoc \
	build/cli/regdoctype

build/lambda/% : cmd/lambda/*/%.go
	go build -o $@ $?
	zip $@-$(today).zip $@

lambda: build/lambda/admin \
	build/lambda/changes \
	build/lambda/getdoc \
	build/lambda/listen \
	build/lambda/sqs-listen \
	build/lambda/sqs-postdoc

build: lambda cli

all: build

clean:
	rm -r ./build
