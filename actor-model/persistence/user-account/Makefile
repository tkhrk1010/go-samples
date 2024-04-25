start:
	go run main.go

test:
	go test -v ./...

.PHONY: proto
proto:
	cd ./proto && \
	sh ./build.sh && \
	cd ..