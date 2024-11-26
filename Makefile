VERSION := $(shell git describe --tags --always --dirty)
REPO := https://api.github.com/repos/andrewrynhard-audio/bpm/releases/latest
OUTPUT_DIR := ./build/bin
RELEASE_DIR := ./build/release

all: build

.PHONY: build
build:
	@mkdir -p $(OUTPUT_DIR)
	@mkdir -p $(RELEASE_DIR)
	@for OS in darwin windows; do \
		for ARCH in amd64 arm64; do \
			wails build -clean -platform $$OS/$$ARCH -ldflags="-X main.Version=$(VERSION) -X main.Repo=$(REPO)"; \
			if [ $$OS = "windows" ]; then \
				tar -czf $(RELEASE_DIR)/bpm-$$OS-$$ARCH.tar.gz -C $(OUTPUT_DIR) bpm.exe; \
			else \
				tar -czf $(RELEASE_DIR)/bpm-$$OS-$$ARCH.tar.gz -C $(OUTPUT_DIR) bpm.app; \
			fi; \
		done; \
	done

.PHONY: tag
tag:
	@if [ -z "$(TAG)" ]; then \
		echo "TAG is not set"; \
		exit 1; \
	fi
	@git tag -a $(TAG) -m "Release $(TAG)."
	@git push origin $(TAG)

.PHONY: release
release: clean tag build
	@if [ -z "$(GITHUB_TOKEN)" ]; then \
		echo "GITHUB_TOKEN is not set"; \
		exit 1; \
	fi
	@gh release create $(TAG) \
		$(RELEASE_DIR)/* \
		--title "$(TAG)"

.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR)
	rm -rf $(RELEASE_DIR)
	rm -rf frontend/node_modules
	rm -rf frontend/dist
