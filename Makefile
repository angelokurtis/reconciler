.PHONY: help
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

.PHONY: copyright
copyright: get-addlicense ## Ensures source code files have copyright license headers.
	$(ADDLICENSE) -ignore "./vendor/**" -c "" -l "apache" $(shell find -regex '.*\.\(go\|yml\|yaml\|sh\)')

.PHONY: lint
lint: get-golangci-lint ## Ensures code quality through static analysis.
	$(GOLANGCI_LINT) run ./...

ADDLICENSE = $(shell pwd)/bin/addlicense
.PHONY: get-addlicense
get-addlicense: ## Download addlicense locally if necessary.
	$(call go-get-tool,$(ADDLICENSE),github.com/google/addlicense@v1.0.0)

GOLANGCI_LINT = $(shell pwd)/bin/golangci-lint
.PHONY: get-golangci-lint
get-golangci-lint: ## Download golangci-lint locally if necessary.
	$(call go-get-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.2)

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
