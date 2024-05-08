.PHONY: release build test
release:
	./scripts/release.sh

build:
	cd dist && go build .. && cd ..

test:
	go test -v ./...
