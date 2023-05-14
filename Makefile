.PHONY: test
test:
	go test -tags=assert -race ./...

.PHONY: lint
lint:
	golangci-lint run

# clone plugin-pb repo
.PHONY: clone
clone:
	git clone https://github.com/cloudquery/plugin-pb

.PHONY: gen-proto
gen-proto:
	protoc --proto_path=. --go_out . --go_opt=module="github.com/cloudquery/plugin-pb-go" --go-grpc_out=. --go-grpc_opt=module="github.com/cloudquery/plugin-pb-go" plugin-pb/discovery/v0/discovery.proto
	protoc --proto_path=. --go_out . --go_opt=module="github.com/cloudquery/plugin-pb-go" --go-grpc_out=. --go-grpc_opt=module="github.com/cloudquery/plugin-pb-go" plugin-pb/base/v0/base.proto plugin-pb/destination/v0/destination.proto plugin-pb/source/v0/source.proto
	protoc --proto_path=. --go_out . --go_opt=module="github.com/cloudquery/plugin-pb-go" --go-grpc_out=. --go-grpc_opt=module="github.com/cloudquery/plugin-pb-go" plugin-pb/source/v1/source.proto
	protoc --proto_path=. --go_out . --go_opt=module="github.com/cloudquery/plugin-pb-go" --go-grpc_out=. --go-grpc_opt=module="github.com/cloudquery/plugin-pb-go" plugin-pb/source/v2/source.proto
	protoc --proto_path=. --go_out . --go_opt=module="github.com/cloudquery/plugin-pb-go" --go-grpc_out=. --go-grpc_opt=module="github.com/cloudquery/plugin-pb-go" plugin-pb/destination/v1/destination.proto