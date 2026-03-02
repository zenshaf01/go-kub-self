# check to see if we can use ash, in alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

run:
	# The below pipes the first program's output towards stdOut to second programs StdIn
	go run apis/services/sales/main.go | go run apis/tooling/logfmt/main.go

help:
	go run apis/services/sales/main.go --help

version:
	go run apis/services/sales/main.go --version


tidy:
	go mod tidy
	# This is putting all third part code packages in the vendor folder.
	go mod vendor


# Create a K8s cluster with Kind

GOLANG          := golang:1.22
ALPINE          := alpine:3.19
KIND            := kindest/node:v1.29.2
POSTGRES        := postgres:16.2
GRAFANA         := grafana/grafana:10.4.0
PROMETHEUS      := prom/prometheus:v2.51.0
TEMPO           := grafana/tempo:2.4.0
LOKI            := grafana/loki:2.9.0
PROMTAIL        := grafana/promtail:2.9.0

KIND_CLUSTER    := ardan-starter-cluster
NAMESPACE       := sales-system
SALES_APP       := sales
AUTH_APP        := auth
BASE_IMAGE_NAME := localhost/ardanlabs
VERSION         := 0.0.1
SALES_IMAGE     := $(BASE_IMAGE_NAME)/$(SALES_APP):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/metrics:$(VERSION)
AUTH_IMAGE      := $(BASE_IMAGE_NAME)/$(AUTH_APP):$(VERSION)

dev-up:
	# The kind command needs the kind image (docker image for kind), the name of the cluster and the config file for dev
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml # The kind-config is used in dev to open ports in the machine. Never do this in prod.
		# everytime you change the yaml file you need to shutdown the cluster and then bring it back up again.

	# we wait for the above complete
	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

dev-down:
	# Kills everything in the cluster nad brings it down.
	kind delete cluster --name $(KIND_CLUSTER)

# This is used for monitoring the cluster
dev-status-all:
	kubectl get nodes -o wide # list all nodes. `-o wide` provides extra info
	kubectl get svc -o wide # list all services in the current K8s namespace `-o wide` provides extra info
	kubectl get pods -o wide --watch --all-namespaces # list all pods for all namespaces `-o wide` provides extra info

# Monitors the cluster
dev-status:
	# runs every 2 seconds
	watch -n 2 kubectl get pods -o wide --all-namespaces

# Build Docker Image for a service
build: sales

sales:
	# build docker image
	docker build \
		-f zarf/docker/Dockerfile.sales \
		-t $(SALES_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

# Load images into local repo using kind
dev-load:
	# this stages the image in a local docker repo and tells kubernetes to fetch image from here instead of downloading from the net.
	kind load docker-image $(SALES_IMAGE) --name $(KIND_CLUSTER)

# Apply config to k8s using kubectl
# run this when there are config changes for k8s
dev-apply:
	# Kustomize builds a single yaml file which we can feed to kubernetes to deploy our pods and then we pipe into kubectl to apply
	kustomize build zarf/k8s/dev/sales | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(SALES_APP) --timeout=120s --for=condition=Ready

# This restarts the deployment
dev-restart:
	kubectl rollout restart deployment $(SALES_APP) --namespace=$(NAMESPACE)

# rebuild the image binary, load the docker image and restart the pods. (For code updates)
dev-update: build dev-load dev-restart

# build the binary, load the image and then apply it to the cluster (For yaml updates)
dev-update-apply: build dev-load dev-apply

# get logs for deployment
dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(SALES_APP) --all-containers=true -f --tail=100 --max-log-requests=6

# describe the deployment in a namespace
dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(SALES_APP)

# describe the pod that has your app running in.
# If something isnt working correctly, we should run this to see whats actually going on.
dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(SALES_APP)

metrics:
	expvarmon -ports="localhost:3010" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

statsviz:
	open -a "Google Chrome" http://localhost:3010/debug/statsviz