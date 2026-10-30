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
	protoc --proto_path=. --go_out . --go_opt=module="github.com/cloudquery/plugin-pb-go" --go-grpc_out=. --go-grpc_opt=module="github.com/cloudquery/plugin-pb-go" plugin-pb/discovery/v1/discovery.proto
	protoc --proto_path=. --go_out . --go_opt=module="github.com/cloudquery/plugin-pb-go" --go-grpc_out=. --go-grpc_opt=module="github.com/cloudquery/plugin-pb-go" plugin-pb/plugin/v3/plugin.proto
