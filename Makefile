test:
	go test -coverprofile cover.out ./cmd/... && go tool cover -html=cover.out

build:
	go build .