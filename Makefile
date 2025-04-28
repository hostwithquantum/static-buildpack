.PHONY: build clean package setup test test-hugo test-mdbook

builder:=paketobuildpacks/builder-jammy-base:latest
bp:=static-buildpack

pack_cmd:=pack -v

build:
	GOOS=linux GOARCH=amd64 goreleaser build --single-target --clean --snapshot
	cp dist/build_linux_amd64_v1/build ./bin/build
	cp dist/detect_linux_amd64_v1/detect ./bin/detect

clean:
	rm -rf dist/
	rm -f bin/build
	rm -f bin/detect
	pack buildpack remove $(bp) || true

setup:
	pack config default-builder $(builder)
	pack config trusted-builders add $(builder)

package: setup build
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
		--buildpack ./meta-buildpack

test-mdbook:
	$(pack_cmd) build \
		test-mdbook-app \
		--builder $(builder) \
		--path ./tests/mdbook-example \
		-e BP_LOG_LEVEL=DEBUG \
		--buildpack ./meta-buildpack

test: setup test-hugo test-mdbook
