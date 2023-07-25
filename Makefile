PHONY: generate
generate:
	mkdir -p proto
	protoc --go_out=proto --go_opt=paths=import \
	--go-grpc_out=proto --go-grpc_opt=paths=import \
	api/thumbnails.proto
	mv proto/github.com/fluxx1on/go_cli/proto/* proto/
	rm -rf proto/github.com

build:
	go build -C . -o ./client

run-sync:
	make build
	./client --urls="https://www.youtube.com/watch?v=NeQM1c-XCDc&ab_channel=RammsteinOfficial,https://www.youtube.com/watch?v=t023ryQgguQ&ab_channel=GOBELINSParis"

run-async:
	make build
	./client --async --urls="https://www.youtube.com/watch?v=NeQM1c-XCDc&ab_channel=RammsteinOfficial,https://www.youtube.com/watch?v=t023ryQgguQ&ab_channel=GOBELINSParis"