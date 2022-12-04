AUTO_INSTRUMENTATION_JAVA_VERSION ?= "$(shell grep -v '\#' versions.txt | grep autoinstrumentation-java | awk -F= '{print $$2}')"
AUTO_INSTRUMENTATION_PYTHON_VERSION ?= "$(shell grep -v '\#' versions.txt | grep autoinstrumentation-python | awk -F= '{print $$2}')"
AUTO_INSTRUMENTATION_DOTNET_VERSION ?= "$(shell grep -v '\#' versions.txt | grep autoinstrumentation-dotnet | awk -F= '{print $$2}')"

VERSION ?= "$(shell git describe --tags | sed 's/^v//')"
VERSION_DATE ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
VERSION_PKG ?= "github.com/open-telemetry/opentelemetry-operator/internal/version"
TARGETALLOCATOR_VERSION ?= "$(shell grep -v '\#' versions.txt | grep targetallocator | awk -F= '{print $$2}')"

# Image URL to use all building/pushing image targets
IMG_PREFIX ?= sprkldev
IMG_REPO ?= k8s-operator
IMG ?= ${IMG_PREFIX}/${IMG_REPO}
LATEST_TAG ?= ${IMG}:latest
VERSION_TAG ?= ${IMG}:${VERSION}

ARCH ?= $(shell go env GOARCH)

build:
	docker buildx build --platform linux/${ARCH} -t ${LATEST_TAG} -t ${VERSION_TAG} --build-arg VERSION_PKG=${VERSION_PKG} --build-arg VERSION=${VERSION} --build-arg VERSION_DATE=${VERSION_DATE} --build-arg OTELCOL_VERSION=${OTELCOL_VERSION} --build-arg TARGETALLOCATOR_VERSION=${TARGETALLOCATOR_VERSION} --build-arg AUTO_INSTRUMENTATION_JAVA_VERSION=${AUTO_INSTRUMENTATION_JAVA_VERSION} --build-arg AUTO_INSTRUMENTATION_PYTHON_VERSION=${AUTO_INSTRUMENTATION_PYTHON_VERSION} --build-arg AUTO_INSTRUMENTATION_DOTNET_VERSION=${AUTO_INSTRUMENTATION_DOTNET_VERSION} .

push:
	docker push ${VERSION_TAG}
	docker push ${LATEST_TAG}
