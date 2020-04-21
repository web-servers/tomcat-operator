#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ARCH=amd64
UNAMEARCH=`uname -p`
if [ "$UNAMEARCH" == "aarch64" ]; then
  ARCH=arm64
fi
echo "using ARCH: ${ARCH}"

GOOS=${1-linux}
BIN_DIR=${2-$(pwd)/build/_output/bin}
mkdir -p ${BIN_DIR}
PROJECT_NAME="tomcat-operator"
REPO_PATH="github.com/tomcat-operator"
BUILD_PATH="${REPO_PATH}/cmd/manager"
VERSION="$(git describe --tags --always --dirty)"
GO_LDFLAGS="-X ${REPO_PATH}/version.Version=${VERSION}"
echo "building ${PROJECT_NAME}..."
GOOS=${GOOS} GOARCH=${ARCH} CGO_ENABLED=0 go build -o ${BIN_DIR}/${PROJECT_NAME} -ldflags "${GO_LDFLAGS}" $BUILD_PATH
