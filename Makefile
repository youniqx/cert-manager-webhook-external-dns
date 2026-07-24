# The version spec of the Kubernetes control-plane binaries used by the
# integration tests. Should track the k8s.io/* library versions in go.mod.
# setup-envtest itself is tracked as a `tool` dependency in go.mod.
ENVTEST_K8S_VERSION ?= 1.36.x

CHART_DIR := deploy/charts/external-dns-webhook

all: lint test build ## Run all targets

.PHONY: build
build: docker ## Run all builds

.PHONY: lint
lint: vet lint-chart ## Run all formatters and linters

.PHONY: test
test: test-go test-chart ## Run the all tests

.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint-chart
lint-chart: ## Run ct lint
	ct lint --all --config=deploy/config.yaml

.PHONY: test-go
test-go: ## Run the Go test suite (envtest binaries are downloaded automatically)
	export KUBEBUILDER_ASSETS="$$(go tool setup-envtest use $(ENVTEST_K8S_VERSION) -p path)" && \
	export TEST_ASSET_ETCD="$$KUBEBUILDER_ASSETS/etcd" && \
	export TEST_ASSET_KUBE_APISERVER="$$KUBEBUILDER_ASSETS/kube-apiserver" && \
	export TEST_ASSET_KUBECTL="$$KUBEBUILDER_ASSETS/kubectl" && \
	go test $(GOTESTFLAGS) ./...

.PHONY: test-chart
test-chart: ## Run the Helm chart unit tests (requires the helm-unittest plugin)
	helm unittest $(CHART_DIR) --strict

.PHONY: docker
docker: ## Build the multi-platform container image
	docker buildx build --platform linux/amd64,linux/arm64 .
