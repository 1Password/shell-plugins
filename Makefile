new-plugin:
	go run cmd/new/main.go plugin

%/example-secrets:
	go run cmd/new/main.go $@

%/validate:
	go run cmd/new/main.go $@

test:
	go test ./...