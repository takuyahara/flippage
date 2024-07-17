.SILENT:
.PHONY: build-macos
## Build app for macOS
build-macos:
	rm -rf deploy && mkdir deploy
	GOOS=darwin GOARCH=arm64 go build -v -o deploy/flippage

BRANCH := $(shell git symbolic-ref --short HEAD)
VERSION := $(shell echo $(BRANCH) | sed 's/release\///g')
.SILENT:
.PHONY: _is-release-branch
# Verify current branch name
_is-release-branch:
	if [ "$(shell echo $(BRANCH) | grep -E "^release/[0-9]+\.[0-9]+\.[0-9]+$$")" ]; then \
		exit 0; \
	fi; \
	echo "Invalid branch name. Aborted."; \
	exit 1

.SILENT:
.PHONY: _version-readme
# Update version in README.md
_version-readme: _is-release-branch
	sed -i -E 's/(https\:\/\/img\.shields\.io\/badge\/version\-)([0-9]+\.[0-9]+\.[0-9]+)(.*)/\1$(VERSION)\3/g' ./README.md
	printf "Updated version to \033[36m$(VERSION)\033[0m in README.md.\n"


.SILENT:
.PHONY: version
## Update version in files based on branch name
version: _version-readme

.DEFAULT_GOAL := help
.SILENT:
.PHONY: help
help:
	@echo "$$(tput setaf 2)Available rules:$$(tput sgr0)";sed -ne"/^## /{h;s/.*//;:d" -e"H;n;s/^## /---/;td" -e"s/:.*//;G;s/\\n## /===/;s/\\n//g;p;}" ${MAKEFILE_LIST}|awk -F === -v n=$$(tput cols) -v i=4 -v a="$$(tput setaf 6)" -v z="$$(tput sgr0)" '{printf"- %s%s%s\n",a,$$1,z;m=split($$2,w,"---");l=n-i;for(j=1;j<=m;j++){l-=length(w[j])+1;if(l<= 0){l=n-i-length(w[j])-1;}printf"%*s%s\n",-i," ",w[j];}}'