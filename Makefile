# These should be installed
GOCMD:=go
NPMCMD:=npm
DOCKERCMD:=docker
# production migration
PGDUMP:=pg_dump
RSYNCCMD:=rsync
TERN:=${GOPATH}/bin/tern

# Alias some go commands
GOBUILD:=$(GOCMD) build
GOCLEAN:=$(GOCMD) clean
GOTEST:=$(GOCMD) test
GOGET:=$(GOCMD) get
GOMODDOWNLOAD:=$(GOCMD) mod download

# Name of the remote rsync production host for migrations
BASTION=bastion

# Fixed names
APP_DIR:=cmd/nrtm4client
BINARY_NAME_APP:=nrtm4client
BINARY_NAME_APP_UNIX:=$(BINARY_NAME_APP)_unix
BINARY_NAME_DEBUG:=__debug_*

# Image release params
DOCKERFILE_APP_DIR:=deployments/docker/nrtm4client
IMAGE_APP_NAME_RELEASE:=eu.gcr.io/fourth-flag-253822/nrtm4client
IMAGE_APP_NAME:=nrtm4client-dev

# Util
CHECK_VCS:=scripts/checkvcs.sh

MAKEFLAGS += --silent

.PHONY: build build-linux buildgo buildweb checkvcs clean cleanall deploy docker-app-prep coverage emptydb install list migrate migrate-production release rewinddb run test testgo testweb testimage webdev

defaulttarget: build

list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

web/node_modules: ; cd web && $(NPMCMD) install

buildgo:
	$(GOMODDOWNLOAD)
	cd $(APP_DIR) && $(GOBUILD) -o $(BINARY_NAME_APP) -v

build: buildgo

build-linux:
	$(GOMODDOWNLOAD) && \
	cd $(APP_DIR) && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	$(GOBUILD) -o $(BINARY_NAME_APP_UNIX) -v

dumpdb: ; $(PGDUMP) -h localhost --data-only nrtm4 | gzip | > work/nrtm4-data-$(shell date -I).psql.gz

emptydb: ; $(TERN) migrate --destination 1 --config third_party/tern/tern.conf --migrations third_party/tern

migrate: ; $(TERN) migrate --config third_party/tern/tern.conf --migrations third_party/tern

rewinddb: ;	$(TERN) migrate --destination -1 --config third_party/tern/tern.conf --migrations third_party/tern

emptytestdb: ; $(TERN) migrate --destination 0 --config third_party/tern/tern.test.conf --migrations third_party/tern

migratetest: ; $(TERN) migrate --config third_party/tern/tern.test.conf --migrations third_party/tern

# docker-app-prep: buildweb build-linux
# 	mkdir -p $(DOCKERFILE_WEB_DIR)/app
# 	$(RSYNCCMD) -q -a --delete web/www $(DOCKERFILE_WEB_DIR)/app
# 	cp $(CD_WEB_DIR)/$(BINARY_NAME_WEB_UNIX) $(DOCKERFILE_WEB_DIR)/app
# 	mkdir -p $(DOCKERFILE_API_DIR)/app
# 	cp $(APP_DIR)/$(BINARY_NAME_API_UNIX) $(DOCKERFILE_API_DIR)/app

coverage: ;	sh scripts/coverage.sh

testgo: emptytestdb migratetest buildgo
	$(GOTEST) ./internal/...

testweb: web/node_modules
	cd web && $(NPMCMD) run test

test: testgo

install: test docker-app-prep
	cd $(DOCKERFILE_APP_DIR) && $(DOCKERCMD) build -t $(IMAGE_APP_NAME) .

testimage: install
	-$(DOCKERCMD) stop nrtm4client_testcontainer 2>/dev/null
	cd $(DOCKERFILE_WEB_DIR) && $(DOCKERCMD) run -dp 8000:8080 --rm --name nrtm4client_testcontainer --env-file ./env.conf $(IMAGE_WEB_NAME) >/dev/null
	#cd web && $(NPMCMD) run e2e >/dev/null
	$(DOCKERCMD) stop nrtm4client_testcontainer >/dev/null

checkvcs:
	sh scripts/checkvcs.sh || (echo "You have local changes to your files. Synchronize your changes and try again."; exit 1)

release-app: checkvcs install
	$(DOCKERCMD) tag $(IMAGE_APP_NAME) $(IMAGE_APP_NAME_RELEASE):$(shell git rev-parse --short HEAD)
	$(DOCKERCMD) push $(IMAGE_APP_NAME_RELEASE):$(shell git rev-parse --short HEAD)

release: release-app ;

clean:
	$(GOCLEAN) ./...
#	rm -rf web/www
#	rm -rf $(DOCKERFILE_APP_DIR)/app
	rm -f $(APP_DIR)/$(BINARY_NAME_DEBUG) $(APP_DIR)/$(BINARY_NAME_APP) $(APP_DIR)/$(BINARY_NAME_APP_UNIX)
#	-$(DOCKERCMD) image rm $(IMAGE_APP_NAME) >/dev/null 2>&1
#	-$(DOCKERCMD) rmi $(shell docker images --filter=reference="$(IMAGE_APP_NAME_RELEASE):*" -q) 2>/dev/null

cleanall: clean ;

uninstall: cleanall ;
