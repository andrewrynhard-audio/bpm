VERSION := $(shell git describe --tags --always --dirty)
REPO := https://api.github.com/repos/andrewrynhard-audio/bpm/releases/latest
OUTPUT_DIR := ./build/bin
RELEASE_DIR := ./build/release

define guard
	@if [ -z "$($(1))" ]; then \
		echo "Error: $(1) is not set"; \
		exit 1; \
	fi
endef

check-env:
	$(call guard,TAG)
	$(call guard,GITHUB_TOKEN)

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
	@git tag -a $(TAG) -m "Release $(TAG)."
	@git push origin $(TAG)

.PHONY: release
release: check-env clean tag build
	@gh release create $(TAG) \
		$(RELEASE_DIR)/* \
		--title "$(TAG)"

.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR)
	rm -rf $(RELEASE_DIR)
	rm -rf frontend/node_modules
	rm -rf frontend/dist
