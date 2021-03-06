OUTPUT=./releases
REPO=github.com/jysperm/deploybeta

export WORKDIR=$(shell pwd)

all: pack-tar

binaries:
	GOOS=linux go build -o $(OUTPUT)/deploybeta-linux-amd64

resources:
	cd frontend && gulp
	mkdir -p $(OUTPUT)/frontend/public
	cp -r frontend/public $(OUTPUT)/frontend
	rm -r $(OUTPUT)/frontend/public/scripts
	cp -r assets $(OUTPUT)

pack-tar: binaries resources
	cd $(OUTPUT) && tar --exclude *.tar.gz -zcvf deploybeta-linux-amd64.tar.gz *

test:
	go test -v $(REPO)
	go test -v $(REPO)/config
	go test -v $(REPO)/lib/builder
	go test -v $(REPO)/lib/datasources
	go test -v $(REPO)/lib/db
	go test -v $(REPO)/lib/models
	go test -v $(REPO)/lib/runtimes
	go test -v $(REPO)/lib/swarm
	go test -v $(REPO)/lib/testing
	go test -v $(REPO)/lib/utils
	go test -v $(REPO)/web/handlers
	go test -v $(REPO)/web/handlers/helpers
	go test -v $(REPO)/web/tests

clean:
	rm -r releases
