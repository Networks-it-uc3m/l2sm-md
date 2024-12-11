# Image URL to use for building/pushing
IMG ?= alexdecb/l2sm-md:0.2
# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.29.0

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

REPOSITORY=l2sm-md
## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUBECTL ?= kubectl
KUSTOMIZE ?= $(LOCALBIN)/kustomize-$(KUSTOMIZE_VERSION)
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen-$(CONTROLLER_TOOLS_VERSION)
ENVTEST ?= $(LOCALBIN)/setup-envtest-$(ENVTEST_VERSION)
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
KIND ?= kind
DOCKER ?= docker 

WORKER_CLUSTER_NUM ?= 2
## Tool Versions
KUSTOMIZE_VERSION ?= v5.3.0
CONTROLLER_TOOLS_VERSION ?= v0.14.0
ENVTEST_VERSION ?= latest
GOLANGCI_LINT_VERSION ?= v1.54.2
##@ Build and Push


L2SMMD_NAMESPACE ?= default

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
export PATH := $(PATH):$(LOCALBIN)
generate-proto: install-tools ## Generate gRPC code from .proto file.
	protoc -I=api/v1 --go_out=paths=source_relative:./api/v1/l2smmd --go-grpc_out=paths=source_relative:./api/v1/l2smmd api/v1/l2smmd.proto

.PHONY: run
include .env
export $(shell sed 's/=.*//' .env)
run: 
	go run ./cmd/server

.PHONY: build
build: fmt vet 
	go build -o $(LOCALBIN)/server ./cmd/server/
	go build -o $(LOCALBIN)/apply-cert ./cmd/apply-cert/

.PHONY: build-installer
build-installer: kustomize ## Generate a consolidated YAML with CRDs and deployment.
	echo "" > deployments/l2smmd-deployment.yaml
	echo "---" >> deployments/l2smmd-deployment.yaml  # Add a document separator before appending
	cd config/server && $(KUSTOMIZE) edit set image server=${IMG}
	$(KUSTOMIZE) build config/default >> deployments/l2smmd-deployment.yaml

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...


.PHONY: deploy
deploy: kustomize ## Deploy server to the K8s cluster specified in ~/.kube/config.
	cd config/server && $(KUSTOMIZE) edit set image server=${IMG}
	$(KUSTOMIZE) build config/default | $(KUBECTL) apply -f - 

.PHONY: undeploy
undeploy: kustomize ## Undeploy server from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/default | $(KUBECTL) delete --ignore-not-found=true -f -

.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	$(call go-install-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v5,$(KUSTOMIZE_VERSION))

.PHONY: setup-dev
setup-dev: create-control-plane create-workers add-cni install-l2sm 
	$(KUBECTL) config use-context kind-control-plane

.PHONY: deploy-dev
deploy-dev: apply-cert kustomize
	$(KUSTOMIZE) build config/dev | $(KUBECTL) apply -f - 
	
.PHONY: undeploy-dev
undeploy-dev: kustomize ## Undeploy server from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/dev | $(KUBECTL) delete --ignore-not-found=true -f -
	$(KUBECTL) delete secrets --all

# Define file extensions for various formats
FILES := $(shell find . -type f \( -name "*.go" -o -name "*.json" -o -name "*.yaml" -o -name "*.yml" -o -name "*.md" \))

# Install the addlicense tool if not installed
.PHONY: install-tools
install-tools:
	GOBIN=$(LOCALBIN) go install github.com/google/addlicense@latest
	GOBIN=$(LOCALBIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCALBIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


.PHONY: add-license
add-license: install-tools
	@for file in $(FILES); do \
		$(LOCALBIN)/addlicense -f ./hack/LICENSE.txt -l apache "$${file}"; \
	done


CERT_FILES := $(shell find ./config/certs/ -name "*.key")

.PHONY: apply-cert
apply-cert: build
	@if [ -n "$(CERT_FILES)" ]; then \
		for file in $(CERT_FILES); do \
			clustername=$$(basename $${file} .key); \
			$(LOCALBIN)/apply-cert --namespace $(L2SMMD_NAMESPACE) --clustername $${clustername} $${file}; \
		done; \
	else \
		echo "No certificate files to process."; \
	fi

.PHONY: create-control-plane
create-control-plane:
	$(KIND) create cluster --config ./examples/quickstart/control-plane-cluster.yaml

.PHONY: create-workers
create-workers:
	for number in $(shell seq 1 ${WORKER_CLUSTER_NUM}); do \
		$(KIND) create cluster --config ./examples/quickstart/worker-cluster.yaml --name worker-cluster-$$number; \
		$(KUBECTL) config view -o jsonpath='{.clusters[?(@.name == "kind-worker-cluster-'$$number'")].cluster.certificate-authority-data}' --raw | base64 -d > config/certs/kind-worker-cluster-$$number.key; \
	done

.PHONY: install-l2sm
install-l2sm:
	for number in $(shell seq 1 ${WORKER_CLUSTER_NUM}); do \
		$(KUBECTL) apply --context kind-worker-cluster-$$number -f https://github.com/cert-manager/cert-manager/releases/download/v1.16.1/cert-manager.yaml; \
		$(KUBECTL) apply --context kind-worker-cluster-$$number -f https://raw.githubusercontent.com/k8snetworkplumbingwg/multus-cni/master/deployments/multus-daemonset-thick.yml; \
	done; \
	for number in $(shell seq 1 ${WORKER_CLUSTER_NUM}); do \
		$(KUBECTL) --context kind-worker-cluster-$$number wait --for=condition=Ready pods --all -A --timeout=300s; \
	done; \
	for number in $(shell seq 1 ${WORKER_CLUSTER_NUM}); do \
		$(KUBECTL) --context kind-worker-cluster-$$number create -f https://github.com/Networks-it-uc3m/L2S-M/raw/refs/heads/main/deployments/l2sm-deployment.yaml; \
	done

.PHONY: add-cni
add-cni:
	@if [ ! -d "plugins/bin" ] || [ -z "$$(ls -A plugins/bin)" ]; then \
		mkdir -p plugins/bin; \
		wget -q https://github.com/containernetworking/plugins/releases/download/v1.6.0/cni-plugins-linux-amd64-v1.6.0.tgz; \
		tar -xf cni-plugins-linux-amd64-v1.6.0.tgz -C plugins/bin; \
		rm cni-plugins-linux-amd64-v1.6.0.tgz; \
	fi
	@nodes="$$( $(KIND) get nodes -A)"; \
	if [ -z "$$nodes" ]; then \
		echo "No nodes found. Is Kind running?"; \
		exit 1; \
	fi; \
	for node in $$nodes; do \
		echo "Copying plugins to node: $$node"; \
		$(DOCKER) cp ./plugins/bin/. $$node:/opt/cni/bin; \
		if [ $$? -ne 0 ]; then \
			echo "Failed to copy plugins to $$node"; \
			exit 1; \
		fi; \
		$(DOCKER) exec $$node modprobe br_netfilter; \
		$(DOCKER) exec $$node sysctl -p /etc/sysctl.conf; \
	done; \
	clusters="$$( $(KIND) get clusters )"; \
	if [ -z "$$clusters" ]; then \
		echo "No clusters found. Is Kind running?"; \
		exit 1; \
	fi; \
	for cluster in $$clusters; do \
		$(KUBECTL) --context kind-$$cluster apply -f https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml; \
		if [ $$? -ne 0 ]; then \
			echo "Failed to install flannel in $$cluster"; \
			exit 1; \
		fi; \
	done

.PHONY: copy-to-container
copy-to-container:
	@if [ -z "$(container)" ]; then \
		echo "Error: Please specify a container name using 'make copy-to-container container=<container_name>'"; \
		exit 1; \
	fi
	docker cp ./plugins/bin/. $(container):/opt/cni/bin
# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary (ideally with version)
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f $(1) ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv "$$(echo "$(1)" | sed "s/-$(3)$$//")" $(1) ;\
}
endef


.PHONY: clean
clean:
	$(KIND) delete clusters --all
