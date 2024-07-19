SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Release

.PHONY: current-version
current-version: svu ## Show the current version of the project.
	@current_version=$$($(SVU) current); \
	echo $$current_version

.PHONY: next-version
next-version: svu ## Calculate the next version, following semantic versioning.
	@current_version=$$($(SVU) current); \
	next_version=$$($(SVU) next); \
	if [ "$$current_version" = "$$next_version" ]; then \
		echo "Error: Current version is equal to next version."; \
		exit 1; \
	fi; \
	echo $$next_version

.PHONY: auto-release
auto-release: svu ## Automate release tasks based on git log.
	@current_version=$$($(SVU) current); \
	next_version=$$($(SVU) next); \
	if [ "$$current_version" = "$$next_version" ]; then \
		echo "Error: Current version is equal to next version."; \
		exit 1; \
	fi; \
	current_branch=$$(git rev-parse --abbrev-ref HEAD); \
	remote_branch=$$(git for-each-ref --format='%(upstream:short)' refs/heads/"$$current_branch"); \
	if [ "$$remote_branch" != "origin/main" ]; then \
		echo "Error: You are not in the main branch."; \
		exit 1; \
	fi; \
	next_version=$$($(SVU) next); \
	git tag -a $$next_version -m "release $$next_version"; \
	git push origin $$next_version; \
	echo "$$next_version has been released successfully"

.PHONY: major-release
major-release: svu ## Force a major release with significant changes.
	@current_branch=$$(git rev-parse --abbrev-ref HEAD); \
	remote_branch=$$(git for-each-ref --format='%(upstream:short)' refs/heads/"$$current_branch"); \
	if [ "$$remote_branch" != "origin/main" ]; then \
		echo "Error: You are not in the main branch."; \
		exit 1; \
	fi; \
	next_version=$$($(SVU) major); \
	git tag -a $$next_version -m "release $$next_version"; \
	git push origin $$next_version; \
	echo "$$next_version has been released successfully"

.PHONY: minor-release
minor-release: svu ## Force a minor release with new features.
	@current_branch=$$(git rev-parse --abbrev-ref HEAD); \
	remote_branch=$$(git for-each-ref --format='%(upstream:short)' refs/heads/"$$current_branch"); \
	if [ "$$remote_branch" != "origin/main" ]; then \
		echo "Error: You are not in the main branch."; \
		exit 1; \
	fi; \
	next_version=$$($(SVU) minor); \
	git tag -a $$next_version -m "release $$next_version"; \
	git push origin $$next_version; \
	echo "$$next_version has been released successfully"

.PHONY: patch-release
patch-release: svu ## Force a patch release with bug fixes.
	@current_branch=$$(git rev-parse --abbrev-ref HEAD); \
	remote_branch=$$(git for-each-ref --format='%(upstream:short)' refs/heads/"$$current_branch"); \
	if [ "$$remote_branch" != "origin/main" ]; then \
		echo "Error: You are not in the main branch."; \
		exit 1; \
	fi; \
	next_version=$$($(SVU) patch); \
	git tag -a $$next_version -m "release $$next_version"; \
	git push origin $$next_version; \
	echo "$$next_version has been released successfully"

##@ Tool Binaries

SVU = $(shell pwd)/bin/svu
.PHONY: svu
svu: ## Checks for svu installation and downloads it if not found.
	$(call go-get-tool,$(SVU),github.com/caarlos0/svu@v1.11.0)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
