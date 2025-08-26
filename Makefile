.PHONY: build clean setup test test-hugo test-hugo-npm test-mdbook

builder:=paketobuildpacks/builder-jammy-base:latest
bp:=runway-buildpacks/static-websites

pack_cmd:=pack -v

BUILD_DIR?=./build
VERSION?=dev

build:
	GOOS=linux GOARCH=amd64 goreleaser build --single-target --clean --snapshot
	cp dist/build_linux_amd64_v1/build ./bin/build
	cp dist/detect_linux_amd64_v1/detect ./bin/detect

clean: restore-version
	rm -rf dist/
	rm -f bin/build
	rm -f bin/detect
	pack buildpack yank $(bp) || true

setup:
	pack config default-builder $(builder)
	pack config trusted-builders add $(builder)

test-hugo-go-%: webserver=$*
test-hugo-go-%:
	$(pack_cmd) build \
		test-hugo-go-$(webserver)-app \
		--builder $(builder) \
		--path ./tests/hugo-go \
		-e BP_LOG_LEVEL=DEBUG \
		-e BP_WEB_SERVER=$(webserver) \
		-e BP_WEB_SERVER_ROOT=./ \
		--buildpack ./meta-buildpack
	$(info docker run -it --platform linux/amd64 --rm --env PORT=8666 -p 8666:8666 hugo-go-$(webserver)-app)

test-hugo-npm-%: webserver=$*
test-hugo-npm-%:
	$(pack_cmd) build \
		test-hugo-npm-$(webserver)-app \
		--builder $(builder) \
		--path ./tests/hugo-npm \
		-e BP_LOG_LEVEL=DEBUG \
		-e BP_WEB_SERVER=$(webserver) \
		-e BP_WEB_SERVER_ROOT=./ \
		-e BP_NODE_RUN_SCRIPTS=build \
		-e BP_KEEP_FILES=static/style.css \
		--buildpack ./meta-buildpack
	$(info docker run -it --platform linux/amd64 --rm --env PORT=8666 -p 8666:8666 hugo-npm-$(webserver)-app)

test-hugo-%: webserver=$*
test-hugo-%:
	$(pack_cmd) build \
		test-hugo-$(webserver)-app \
		--builder $(builder) \
		--path ./tests/hugo-example \
		-e BP_LOG_LEVEL=DEBUG \
		-e BP_WEB_SERVER=$(webserver) \
		-e BP_WEB_SERVER_ROOT=./ \
		--buildpack ./meta-buildpack
	$(info docker run -it --platform linux/amd64 --rm --env PORT=8666 -p 8666:8666 hugo-$(webserver)-app)

test-mdbook-%: webserver=$*
test-mdbook-%:
	$(pack_cmd) build \
		test-mdbook-$(webserver)-app \
		--builder $(builder) \
		--path ./tests/mdbook-example \
		-e BP_LOG_LEVEL=DEBUG \
		-e BP_WEB_SERVER=$(webserver) \
		-e BP_WEB_SERVER_ROOT=./ \
		--buildpack ./meta-buildpack
	$(info docker run -it --platform linux/amd64 --rm --env PORT=8666 -p 8666:8666 mdbook-$(webserver)-app)

# Legacy targets for backwards compatibility
test-hugo: test-hugo-nginx
test-hugo-npm: test-hugo-npm-nginx
test-mdbook: test-mdbook-nginx

test: package test-hugo-nginx test-hugo-httpd test-hugo-npm-nginx test-hugo-npm-httpd test-mdbook-nginx test-mdbook-httpd

.PHONY: prep
prep:
	mkdir -p $(BUILD_DIR)/bin
	cp dist/build_linux_amd64*/build $(BUILD_DIR)/bin/
	cp dist/detect_linux_amd64*/detect $(BUILD_DIR)/bin/
	cp buildpack.toml $(BUILD_DIR)/
	sed -i.bak -E "s/__replace__/$(VERSION)/" $(BUILD_DIR)/buildpack.toml
	rm -f $(BUILD_DIR)/buildpack.toml.bak
	cp package.toml $(BUILD_DIR)/