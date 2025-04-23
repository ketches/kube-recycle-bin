KRB_VERSION := v0.1.0

ALIYUN_REGISTRY := registry.cn-hangzhou.aliyuncs.com/ketches

.PHONY: install
install:
	@echo "» installing krb-cli..."
	go install -ldflags="-X github.com/ketches/kube-recycle-bin/cmd/krb-cli/cmd.Version=${KRB_VERSION}" ./cmd/krb-cli

.PHONY: build
build: build-binary build-binary build-docker

build-binary: build-controller-binary build-webhook-binary

build-controller-binary:
	@echo "» building krb binary..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64/krb-controller cmd/krb-controller/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/arm64/krb-controller cmd/krb-controller/main.go

build-webhook-binary:
	@echo "» building krb binary..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64/krb-webhook cmd/krb-webhook/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/arm64/krb-webhook cmd/krb-webhook/main.go

docker-buildx-init:
	@echo "» initializing docker buildx..."
	docker buildx create --use --name gobuilder 2>/dev/null || docker buildx use gobuilder

build-docker: docker-buildx-init build-docker-controller build-docker-webhook

build-docker-controller:
	@echo "» building krb-controller docker image..."
	docker buildx build --platform linux/amd64,linux/arm64 \
	--build-arg KRB_APPNAME=krb-controller \
	-t ketches/krb-controller:${KRB_VERSION} \
	-t ketches/krb-controller:latest \
	-t ${ALIYUN_REGISTRY}/krb-controller:${KRB_VERSION} \
	-t ${ALIYUN_REGISTRY}/krb-controller:latest \
	--push . -f Dockerfile.local

build-docker-webhook:
	@echo "» building krb-webhook docker image..."
	docker buildx build --platform linux/amd64,linux/arm64 \
	--build-arg KRB_APPNAME=krb-webhook \
	-t ketches/krb-webhook:${KRB_VERSION} \
	-t ketches/krb-webhook:latest \
	-t ${ALIYUN_REGISTRY}/krb-webhook:${KRB_VERSION} \
	-t ${ALIYUN_REGISTRY}/krb-webhook:latest \
	--push . -f Dockerfile.local

.PHONY: deploy
deploy: deploy-crds
	@echo "» deploying krb controller and webhook..."
	kubectl apply -f manifests/deploy.yaml

deploy-crds:
	@echo "» deploying krb crds..."
	kubectl apply -f manifests/crds.yaml

.PHONY: undeploy
undeploy: undeploy-crds
	@echo "» undeploying krb controller and webhook..."
	kubectl delete -f manifests/deploy.yaml

undeploy-crds:
	@echo "» undeploying krb crds..."
	kubectl delete -f manifests/crds.yaml

.PHONY: release
release:
	@if [ -z "${KRB_VERSION}" ]; then \
		echo "KRB_VERSION is not set"; \
		exit 1; \
	fi
	@if git rev-parse "refs/tags/${KRB_VERSION}" >/dev/null 2>&1; then \
        echo "Git tag ${KRB_VERSION} already exists, please use a new version."; \
        exit 1; \
    fi
	@sed -E -i '' 's/(var Version = ")[^"]+(")/\1${KRB_VERSION}\2/' cmd/krb-cli/cmd/version.go
	@git add cmd/krb-cli/cmd/version.go
	@git commit -m "Release ${KRB_VERSION}"
	@git push
	git tag -a "${KRB_VERSION}" -m "release ${KRB_VERSION}"
	git push origin "${KRB_VERSION}"