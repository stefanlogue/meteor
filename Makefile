.PHONY: release build
release:
	./scripts/release.sh

build:
	cd dist && go build .. && cd ..
