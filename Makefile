# Image URL to use for building/pushing
IMG ?= alexdecb/l2sm-md:0.1

# CONTAINER_TOOL defines the container tool to be used for building images (defaults to docker).
CONTAINER_TOOL ?= docker

# Setting SHELL to bash allows bash commands to be executed by recipes.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

##@ Build and Push

.PHONY: docker-build
docker-build: ## Build docker image with the server.
	$(CONTAINER_TOOL) build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image to the repository.
	$(CONTAINER_TOOL) push ${IMG}

##@ Help

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: generate-proto
generate-proto: ## Generate gRPC code from .proto file.
	protoc -I=api/v1 --go_out=paths=source_relative:./api/v1/l2smmd --go-grpc_out=paths=source_relative:./api/v1/l2smmd api/v1/l2smmd.proto

.PHONY: run
include .env
export $(shell sed 's/=.*//' .env)
run: # manifests generate fmt vet ## Run a controller from your host.
	go run ./cmd/main.go

.PHONY: build
build: manifests generate fmt vet ## Build manager binary.
	go build -o bin/manager cmd/main.go

.PHONY: build-installer
build-installer: kustomize ## Generate a consolidated YAML with CRDs and deployment.
	echo "" > deployments/l2sm-deployment.yaml
	echo "---" >> deployments/l2sm-deployment.yaml  # Add a document separator before appending
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default >> deployments/l2sm-deployment.yaml
	$(KUSTOMIZE) build config/tmp >> deployments/l2sm-deployment.yaml

.PHONY: deploy
deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | $(KUBECTL) apply -f -
	$(KUSTOMIZE) build config/tmp | $(KUBECTL) apply -f -

.PHONY: undeploy
undeploy: kustomize ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/tmp | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -
	$(KUSTOMIZE) build config/default | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	$(call go-install-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v5,$(KUSTOMIZE_VERSION))
