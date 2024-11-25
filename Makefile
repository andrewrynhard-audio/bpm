VERSION := $(shell git describe --tags --always --dirty)
REPO := https://api.github.com/repos/andrewrynhard-audio/bpm/releases/latest

all: build

build:
	go build -ldflags="-X github.com/andrewrynhard-audio/bpm/pkg/ui/update.Version=$(VERSION) -X github.com/andrewrynhard-audio/bpm/pkg/ui/update.Repo=$(REPO)" .

release:
ifeq ($(TAG),)
	$(error TAG is not set)
endif
	git tag -a $(TAG) -m "Release $(TAG)."
	git push origin $(TAG)
	VERSION=$(TAG) REPO=$(REPO) goreleaser release --clean

release-dry-run:
ifeq ($(TAG),)
	$(error TAG is not set)
endif
	VERSION=$(TAG) REPO=$(REPO) goreleaser build --snapshot --clean
