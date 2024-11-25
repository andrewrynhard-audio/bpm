VERSION := $(shell git describe --tags --always --dirty)
REPO := https://api.github.com/repos/andrewrynhard-audio/bpm/releases/latest

build:
	go build -ldflags="-X github.com/andrewrynhard-audio/bpm/pkg/update.Version=$(VERSION) -X github.com/andrewrynhard-audio/bpm/pkg/update.Repo=$(REPO)" .

release:
ifeq ($(TAG),)
	$(error TAG is not set)
endif
	git tag -a $(TAG) -m "Release $(TAG)."
	git push origin $(TAG)
	VERSION=$(TAG) goreleaser release --clean
