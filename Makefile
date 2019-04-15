SHA = $(shell gitmeta git sha)
TAG = $(shell gitmeta image tag)
BUILT = $(shell gitmeta built)
PUSH = $(shell gitmeta pushable)

VPATH = $(PATH)
KERNEL_IMAGE ?= autonomy/kernel:8fd9a83
TOOLCHAIN_IMAGE ?= autonomy/toolchain:989387e
DOCKER_ARGS ?=
BUILDKIT_VERSION ?= v0.4.0
BUILDKIT_IMAGE ?= moby/buildkit:$(BUILDKIT_VERSION)
BUILDKIT_HOST ?= tcp://0.0.0.0:1234
#BUILDKIT_CACHE ?= -v $(HOME)/.buildkit:/var/lib/buildkit
BUILDKIT_CONTAINER_NAME ?= talos-buildkit
BUILDKIT_CONTAINER_STOPPED := $(shell docker ps --filter name=$(BUILDKIT_CONTAINER_NAME) --filter status=exited --format='{{.Names}}' 2>/dev/null)
BUILDKIT_CONTAINER_RUNNING := $(shell docker ps --filter name=$(BUILDKIT_CONTAINER_NAME) --filter status=running --format='{{.Names}}' 2>/dev/null)
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
BUILDCTL_ARCHIVE := https://github.com/moby/buildkit/releases/download/$(BUILDKIT_VERSION)/buildkit-$(BUILDKIT_VERSION).linux-amd64.tar.gz
endif
ifeq ($(UNAME_S),Darwin)
BUILDCTL_ARCHIVE := https://github.com/moby/buildkit/releases/download/$(BUILDKIT_VERSION)/buildkit-$(BUILDKIT_VERSION).darwin-amd64.tar.gz
endif
BINDIR ?= ./bin
CONFORM_VERSION ?= 57c9dbd

COMMON_ARGS = --progress=plain
COMMON_ARGS += --frontend=dockerfile.v0
COMMON_ARGS += --local context=.
COMMON_ARGS += --local dockerfile=.
COMMON_ARGS += --opt build-arg:KERNEL_IMAGE=$(KERNEL_IMAGE)
COMMON_ARGS += --opt build-arg:TOOLCHAIN_IMAGE=$(TOOLCHAIN_IMAGE)
COMMON_ARGS += --opt build-arg:SHA=$(SHA)
COMMON_ARGS += --opt build-arg:TAG=$(TAG)

# to allow tests to run containerd
DOCKER_TEST_ARGS = --security-opt seccomp:unconfined --privileged -v /var/lib/containerd/

all: ci kernel initramfs rootfs osctl-linux-amd64 osctl-darwin-amd64 osinstall-linux-amd64 installer talos

.PHONY: builddeps
builddeps: gitmeta buildctl

gitmeta:
	GO111MODULE=off go get github.com/talos-systems/gitmeta

buildctl: $(BINDIR)/buildctl

$(BINDIR)/buildctl:
	@mkdir $(BINDIR)
	@wget -qO - $(BUILDCTL_ARCHIVE) | \
		tar -zxf - -C $(BINDIR) --strip-components 1 bin/buildctl

.PHONY: buildkitd
buildkitd:
ifeq (tcp://0.0.0.0:1234,$(findstring tcp://0.0.0.0:1234,$(BUILDKIT_HOST)))
ifeq ($(BUILDKIT_CONTAINER_STOPPED),$(BUILDKIT_CONTAINER_NAME))
	@echo "Removing exited talos-buildkit container"
	@docker rm $(BUILDKIT_CONTAINER_NAME)
endif
ifneq ($(BUILDKIT_CONTAINER_RUNNING),$(BUILDKIT_CONTAINER_NAME))
	@echo "Starting talos-buildkit container"
	@docker run \
		--name $(BUILDKIT_CONTAINER_NAME) \
		-d \
		--privileged \
		-p 1234:1234 \
		$(BUILDKIT_CACHE) \
		$(BUILDKIT_IMAGE) \
		--addr $(BUILDKIT_HOST)
	@echo "Wait for buildkitd to become available"
	@sleep 5
endif
endif

enforce:
	@docker run --rm -v $(PWD):/src -w /src autonomy/conform:$(CONFORM_VERSION)

.PHONY: ci
ci: builddeps buildkitd

base: buildkitd
	@$(BINDIR)/buildctl --addr $(BUILDKIT_HOST) \
		build \
		--output type=docker,dest=build/$@.tar,name=docker.io/autonomy/$@:$(TAG) \
		--opt target=$@ \
		$(COMMON_ARGS)
	@docker load < build/$@.tar

kernel: buildkitd
	@$(BINDIR)/buildctl --addr $(BUILDKIT_HOST) \
		build \
    --output type=local,dest=build \
		--opt target=$@ \
		$(COMMON_ARGS)

initramfs: buildkitd
	@$(BINDIR)/buildctl --addr $(BUILDKIT_HOST) \
		build \
    --output type=local,dest=build \
		--opt target=$@ \
		$(COMMON_ARGS)

rootfs: buildkitd hyperkube etcd coredns pause osd trustd proxyd ntpd
	@$(BINDIR)/buildctl --addr $(BUILDKIT_HOST) \
		build \
    --output type=local,dest=build \
		--opt target=$@ \
		$(COMMON_ARGS)

installer: buildkitd
	@mkdir -p build
	@$(BINDIR)/buildctl --addr $(BUILDKIT_HOST) \
		build \
		--output type=docker,dest=build/$@.tar,name=docker.io/autonomy/$@:$(TAG),push=$(PUSH) \
		--opt target=$@ \
		$(COMMON_ARGS)
	@docker load < build/$@.tar

proto: buildkitd
	$(BINDIR)/buildctl --addr $(BUILDKIT_HOST) \
		build \
    --output type=local,dest=./ \
		--opt target=$@ \
		$(COMMON_ARGS)

talos-gce: installer
	@docker run --rm -v /dev:/dev -v $(PWD)/build:/out --privileged $(DOCKER_ARGS) autonomy/installer:$(TAG) disk -l -f -p googlecloud -u none -e 'random.trust_cpu=on'
	@tar -C $(PWD)/build -Sczf $(PWD)/build/$@.tar.gz disk.raw
	@rm $(PWD)/build/disk.raw

talos-raw: installer
	@docker run --rm -v /dev:/dev -v $(PWD)/build:/out --privileged $(DOCKER_ARGS) autonomy/installer:$(TAG) disk -n rootfs -l

talos: buildkitd
	@$(BINDIR)/buildctl --addr $(BUILDKIT_HOST) \
		build \
		--output type=docker,dest=build/$@.tar,name=docker.io/autonomy/$@:$(TAG),push=$(PUSH) \
		--opt target=$@ \
		$(COMMON_ARGS)
	@docker load < build/$@.tar

release: all talos-raw talos-gce talos
	xz -e9 ./build/rootfs.raw

test: buildkitd
	@mkdir -p build
	@$(BINDIR)/buildctl --addr $(BUILDKIT_HOST) \
		build \
		--output type=docker,dest=build/$@.tar,name=docker.io/autonomy/$@:$(TAG),push=$(PUSH) \
		--opt target=$@ \
		$(COMMON_ARGS)
	@docker load < build/$@.tar
	@docker run -i --rm $(DOCKER_TEST_ARGS) autonomy/$@:$(TAG) /bin/test.sh --short
	@trap "rm -rf ./.artifacts" EXIT; mkdir -p ./.artifacts && \
		docker run -i --rm $(DOCKER_TEST_ARGS) -v $(PWD)/.artifacts:/src/artifacts autonomy/$@:$(TAG) /bin/test.sh && \
		cp ./.artifacts/coverage.txt coverage.txt

lint: buildkitd
	@$(BINDIR)/buildctl --addr $(BUILDKIT_HOST) \
		build \
		--opt target=$@ \
		$(COMMON_ARGS)

osctl-linux-amd64: buildkitd
	@$(BINDIR)/buildctl --addr $(BUILDKIT_HOST) \
		build \
    --output type=local,dest=build \
		--opt target=$@ \
		$(COMMON_ARGS)

osctl-darwin-amd64: buildkitd
	@$(BINDIR)/buildctl --addr $(BUILDKIT_HOST) \
		build \
    --output type=local,dest=build \
		--opt target=$@ \
		$(COMMON_ARGS)

osinstall-linux-amd64: buildkitd
	@$(BINDIR)/buildctl --addr $(BUILDKIT_HOST) \
		build \
    --output type=local,dest=build \
		--opt target=$@ \
		$(COMMON_ARGS)

udevd: buildkitd
	@$(BINDIR)/buildctl --addr $(BUILDKIT_HOST) \
		build \
		--opt target=$@ \
		$(COMMON_ARGS)

images:
	mkdir -p images

osd: buildkitd images
	@$(BINDIR)/buildctl --addr $(BUILDKIT_HOST) \
		build \
		--output type=docker,dest=images/$@.tar,name=docker.io/autonomy/$@:$(TAG) \
		--opt target=$@ \
		$(COMMON_ARGS)

trustd: buildkitd images
	@$(BINDIR)/buildctl --addr $(BUILDKIT_HOST) \
		build \
		--output type=docker,dest=images/$@.tar,name=docker.io/autonomy/$@:$(TAG) \
		--opt target=$@ \
		$(COMMON_ARGS)

proxyd: buildkitd images
	@$(BINDIR)/buildctl --addr $(BUILDKIT_HOST) \
		build \
		--output type=docker,dest=images/$@.tar,name=docker.io/autonomy/$@:$(TAG) \
		--opt target=$@ \
		$(COMMON_ARGS)

ntpd: buildkitd images
	@$(BINDIR)/buildctl --addr $(BUILDKIT_HOST) \
		build \
		--output type=docker,dest=images/$@.tar,name=docker.io/autonomy/$@:$(TAG) \
		--opt target=$@ \
		$(COMMON_ARGS)

hyperkube: images
	@docker pull k8s.gcr.io/$@:v1.14.1
	@docker save k8s.gcr.io/$@:v1.14.1 -o ./images/$@.tar

etcd: images
	@docker pull k8s.gcr.io/$@:3.3.10
	@docker save k8s.gcr.io/$@:3.3.10 -o ./images/$@.tar

coredns: images
	@docker pull k8s.gcr.io/$@:1.3.1
	@docker save k8s.gcr.io/$@:1.3.1 -o ./images/$@.tar

pause: images
	@docker pull k8s.gcr.io/$@:3.1
	@docker save k8s.gcr.io/$@:3.1 -o ./images/$@.tar

.PHONY: login
login:
	@docker login --username "$(DOCKER_USERNAME)" --password "$(DOCKER_PASSWORD)"

push:
	@docker tag autonomy/installer:$(TAG) autonomy/installer:latest
	@docker push autonomy/installer:$(TAG)
	@docker push autonomy/installer:latest
	@docker tag autonomy/talos:$(TAG) autonomy/talos:latest
	@docker push autonomy/talos:$(TAG)
	@docker push autonomy/talos:latest

clean:
	@-rm -rf build images vendor
