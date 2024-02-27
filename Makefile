# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL := /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# Image identifier
IMAGE_VERSION ?= latest
DEFAULT_NAMESPACE = quay.io/thvo
IMAGE_NAMESPACE ?= $(DEFAULT_NAMESPACE)

PLATFORM ?= linux/amd64

# Tools
IMAGE_BUILDER ?= podman

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: add-license
add-license: ## Add license header to source file
	$(IMAGE_BUILDER) run \
		--name addlicense \
		--security-opt label=disable \
		-v "$${PWD}":"/src" \
		-it --replace \
		ghcr.io/google/addlicense -v -f ./license.header.txt *.go utils/*.go

##@ Lambda

AWS_REGION ?= ca-central-1

LAMBDA_IMAGE_NAME ?= lambda-auto-deployer
LAMBDA_IMAGE ?= $(IMAGE_NAMESPACE)/$(LAMBDA_IMAGE_NAME):$(IMAGE_VERSION)

.PHONY: oci-build
oci-build: ## Build the lambda OCI image
	$(IMAGE_BUILDER) build -t $(LAMBDA_IMAGE) --platform $(PLATFORM) -f Dockerfile .

.PHONY: oci-push
oci-push: login-ecr ## Push the lambda OCI image.
	AWS_ACCOUNT_ID=$$(aws sts get-caller-identity --query "Account" --output text); \
	ECR_LAMBDA_IMAGE=$${AWS_ACCOUNT_ID}.dkr.ecr.$(AWS_REGION).amazonaws.com/$(LAMBDA_IMAGE_NAME):$(IMAGE_VERSION); \
	$(IMAGE_BUILDER) tag $(LAMBDA_IMAGE) $${ECR_LAMBDA_IMAGE}; \
	$(IMAGE_BUILDER) push $${ECR_LAMBDA_IMAGE}

.PHONY: login-ecr
login-ecr: ## Authenticate container tool with AWS ECR
	AWS_ACCOUNT_ID=$$(aws sts get-caller-identity --query "Account" --output text); \
	aws ecr get-login-password --region $(AWS_REGION) | $(IMAGE_BUILDER) login --username AWS --password-stdin $${AWS_ACCOUNT_ID}.dkr.ecr.$(AWS_REGION).amazonaws.com
