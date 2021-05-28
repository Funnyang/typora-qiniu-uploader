version=v1.1.0
img=""

install:
	@go install -ldflags " \
	-X main.version=$(version) \
    " cmd/tqu.go

run:
	@go run -ldflags " \
	-X main.version=$(version) \
    " cmd/tqu.go '$(img)'