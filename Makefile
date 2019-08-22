DOCKER_REPO ?= docker.io/
IMAGE ?= maxbeck/tomcat-operator
TAG ?= v0.0.1
PROG  := tomcat-operator

.DEFAULT_GOAL := help

## setup            Ensure the operator-sdk is installed.
setup:
	./build/setup-operator-sdk.sh

## tidy             Ensure modules are tidy.
tidy:
	go mod tidy

## codegen          Ensure code is generated.
codegen: setup
	operator-sdk generate k8s
	operator-sdk generate openapi

## build            Compile and build the Tomcat operator.
build: tidy unit-test
	./build/build.sh ${GOOS}

## image            Create the Docker image of the operator
image: build
	docker build -t "${DOCKER_REPO}$(IMAGE):$(TAG)" . -f build/Dockerfile

## push             Push Docker image to the docker.io repository.
push: image
	docker push "${DOCKER_REPO}$(IMAGE):$(TAG)"

## clean            Remove all generated build files.
clean:
	rm -rf build/_output

## run-openshift    Run the Tomcat operator on OpenShift.
run-openshift:
	./build/run-openshift.sh

## run-kubernetes    Run the Tomcat operator on kubernetes.
run-kubernetes:
	./build/run-kubernetes.sh


## test             Perform all tests.
test: unit-test scorecard test-e2e

## scorecard        Run operator-sdk scorecard.
scorecard: setup
	operator-sdk scorecard --verbose

## unit-test        Perform unit tests.
unit-test:
	go test -v ./... -tags=unit

help : Makefile
	@sed -n 's/^##//p' $<
