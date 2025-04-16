.PHONY: build clean package setup test test-hugo test-mdbook
builder?=r.planetary-quantum.com/runway-public/runway-buildpack-stack:jammy-full

platform:=linux/amd64

build:
	GOOS=linux GOARCH=amd64 goreleaser build --single-target --clean --snapshot
	cp dist/build_linux_amd64_v1/build ./bin/build
	cp dist/detect_linux_amd64_v1/detect ./bin/detect

clean:
	rm -rf dist/
	rm -f bin/build
	rm -f bin/detect
	pack buildpack remove static-buildpack || true

setup:
	pack config default-builder $(builder)
	pack config trusted-builders add $(builder)

package: build # DOCKER_DEFAULT_PLATFORM=$(platform)
	 pack -v \
	 	buildpack package \
	 		static-buildpack \
			--config package.toml

test-hugo: setup build
	pack \
		build \
		test-hugo-app \
		--builder $(builder) \
		--platform $(platform) \
		--path ./tests/hugo-example \
		--buildpack .

test-mdbook: setup build
	pack \
		build \
		test-mdbook-app \
		--builder $(builder) \
		--platform $(platform) \
		--path ./tests/mdbook-example \
		--buildpack .

test: test-hugo test-mdbook