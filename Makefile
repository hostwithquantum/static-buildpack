.PHONY: build clean package setup test test-hugo test-hugo-npm test-mdbook prepare-dev restore-version

builder:=paketobuildpacks/builder-jammy-base:latest
bp:=static-buildpack

pack_cmd:=pack -v

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

prepare-dev:
	@if [ ! -f buildpack.toml.bak ]; then \
		cp buildpack.toml buildpack.toml.bak; \
	fi
	@sed -i.tmp -E "s/__replace__/dev/" buildpack.toml && rm buildpack.toml.tmp

restore-version:
	@if [ -f buildpack.toml.bak ]; then \
		mv buildpack.toml.bak buildpack.toml; \
	fi

package: setup build prepare-dev
	$(info packaging $(bp))
	$(pack_cmd) buildpack package \
	 		$(bp) \
			--config package.toml

test-hugo:
	$(pack_cmd) build \
		test-hugo-app \
		--builder $(builder) \
		--path ./tests/hugo-example \
		-e BP_LOG_LEVEL=DEBUG \
		-e BP_WEB_SERVER=nginx \
		-e BP_WEB_SERVER_ROOT=./ \
		--buildpack ./meta-buildpack
	$(info docker run -it --platform linux/amd64 --rm --env PORT=8666 -p 8666:8666 test-hugo-app)

test-hugo-npm:
	$(pack_cmd) build \
		test-hugo-npm-app \
		--builder $(builder) \
		--path ./tests/hugo-npm \
		-e BP_LOG_LEVEL=DEBUG \
		-e BP_WEB_SERVER=nginx \
		-e BP_WEB_SERVER_ROOT=./ \
		-e BP_NODE_RUN_SCRIPTS=build \
		-e BP_KEEP_FILES=static/style.css \
		--buildpack ./meta-buildpack
	$(info docker run -it --platform linux/amd64 --rm --env PORT=8666 -p 8666:8666 test-hugo-npm-app)

test-mdbook:
	$(pack_cmd) build \
		test-mdbook-app \
		--builder $(builder) \
		--path ./tests/mdbook-example \
		-e BP_LOG_LEVEL=DEBUG \
		-e BP_WEB_SERVER=nginx \
		-e BP_WEB_SERVER_ROOT=./ \
		--buildpack ./meta-buildpack
	$(info docker run -it --platform linux/amd64 --rm --env PORT=8666 -p 8666:8666 test-mdbook-app)

test: setup test-hugo test-mdbook
